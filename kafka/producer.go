// Package kafka provides a resilient Kafka producer with circuit breaker pattern.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ProducerConfig holds configuration for the Kafka producer
type ProducerConfig struct {
	Brokers         []string
	ClientID        string
	RequiredAcks    sarama.RequiredAcks
	RetryMax        int
	FlushFrequency  time.Duration
	FlushMessages   int
	CompressionType sarama.CompressionCodec
}

// DefaultProducerConfig returns sensible defaults for banking operations
func DefaultProducerConfig(brokers []string, clientID string) ProducerConfig {
	return ProducerConfig{
		Brokers:         brokers,
		ClientID:        clientID,
		RequiredAcks:    sarama.WaitForAll, // Wait for all replicas (banking requirement)
		RetryMax:        5,
		FlushFrequency:  100 * time.Millisecond,
		FlushMessages:   100,
		CompressionType: sarama.CompressionGZIP,
	}
}

// Producer is a resilient Kafka producer with circuit breaker
type Producer struct {
	producer sarama.SyncProducer
	cb       *gobreaker.CircuitBreaker
	logger   *zap.Logger
	tracer   trace.Tracer
}

// NewProducer creates a new Kafka producer with circuit breaker
func NewProducer(cfg ProducerConfig, logger *zap.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.ClientID = cfg.ClientID
	config.Producer.RequiredAcks = cfg.RequiredAcks
	config.Producer.Retry.Max = cfg.RetryMax
	config.Producer.Flush.Frequency = cfg.FlushFrequency
	config.Producer.Flush.Messages = cfg.FlushMessages
	config.Producer.Compression = cfg.CompressionType
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Enable idempotent producer for exactly-once semantics
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// Configure circuit breaker
	cbSettings := gobreaker.Settings{
		Name:        "kafka-producer",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("Circuit breaker state change",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	}

	return &Producer{
		producer: producer,
		cb:       gobreaker.NewCircuitBreaker(cbSettings),
		logger:   logger,
		tracer:   otel.Tracer("banking-shared/kafka"),
	}, nil
}

// Event interface for publishable events
type Event interface {
	Key() string
}

// Publish sends an event to Kafka with circuit breaker protection
func (p *Producer) Publish(ctx context.Context, topic string, event Event) error {
	ctx, span := p.tracer.Start(ctx, "kafka.publish",
		trace.WithAttributes(
			attribute.String("kafka.topic", topic),
		),
	)
	defer span.End()

	payload, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.Key()),
		Value: sarama.ByteEncoder(payload),
		Headers: []sarama.RecordHeader{
			{Key: []byte("content-type"), Value: []byte("application/json")},
			{Key: []byte("trace-id"), Value: []byte(span.SpanContext().TraceID().String())},
		},
	}

	_, err = p.cb.Execute(func() (interface{}, error) {
		partition, offset, err := p.producer.SendMessage(msg)
		if err != nil {
			return nil, err
		}
		p.logger.Debug("Message sent",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
		return nil, nil
	})

	if err != nil {
		span.RecordError(err)
		p.logger.Error("Failed to publish message",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish to %s: %w", topic, err)
	}

	return nil
}

// PublishBatch sends multiple events to Kafka
func (p *Producer) PublishBatch(ctx context.Context, topic string, events []Event) error {
	ctx, span := p.tracer.Start(ctx, "kafka.publish_batch",
		trace.WithAttributes(
			attribute.String("kafka.topic", topic),
			attribute.Int("batch_size", len(events)),
		),
	)
	defer span.End()

	for _, event := range events {
		if err := p.Publish(ctx, topic, event); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// IsHealthy returns true if the circuit breaker is closed
func (p *Producer) IsHealthy() bool {
	return p.cb.State() == gobreaker.StateClosed
}
