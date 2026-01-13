// Package kafka provides a resilient Kafka consumer with proper error handling.
package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ConsumerConfig holds configuration for the Kafka consumer
type ConsumerConfig struct {
	Brokers  []string
	GroupID  string
	Topics   []string
	ClientID string
}

// MessageHandler is a function that processes a Kafka message
type MessageHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// Consumer is a Kafka consumer group handler
type Consumer struct {
	client  sarama.ConsumerGroup
	handler MessageHandler
	logger  *zap.Logger
	tracer  trace.Tracer
	topics  []string
	ready   chan bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg ConsumerConfig, handler MessageHandler, logger *zap.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = cfg.ClientID
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		client:  client,
		handler: handler,
		logger:  logger,
		tracer:  otel.Tracer("banking-shared/kafka"),
		topics:  cfg.Topics,
		ready:   make(chan bool),
	}, nil
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.client.Consume(ctx, c.topics, c); err != nil {
				c.logger.Error("Consumer error", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready
	c.logger.Info("Consumer started", zap.Strings("topics", c.topics))
	return nil
}

// Stop stops the consumer gracefully
func (c *Consumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	return c.client.Close()
}

// Setup is run at the beginning of a new session
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a partition
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			ctx, span := c.tracer.Start(context.Background(), "kafka.consume",
				trace.WithAttributes(
					attribute.String("kafka.topic", message.Topic),
					attribute.Int64("kafka.partition", int64(message.Partition)),
					attribute.Int64("kafka.offset", message.Offset),
				),
			)

			if err := c.handler(ctx, message); err != nil {
				span.RecordError(err)
				c.logger.Error("Failed to process message",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Error(err),
				)
				// Don't commit on error - message will be reprocessed
				span.End()
				continue
			}

			session.MarkMessage(message, "")
			span.End()

		case <-session.Context().Done():
			return nil
		}
	}
}
