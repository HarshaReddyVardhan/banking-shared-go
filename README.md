# Banking Shared Go Library

Shared Go library for NextGen Banking microservices.

## Installation

Add to your `go.mod`:

```go
require github.com/banking/shared v1.0.0

// For local development, add a replace directive:
replace github.com/banking/shared => ../banking-shared-go
```

## Usage

### Events

```go
import "github.com/banking/shared/events"

// Create a new event
event := events.NewBaseEvent(events.EventTypeTransactionInitiated, "my-service")

// Get topic for event type
topic := events.EventTypeTransactionInitiated.Topic(events.DefaultTopicConfig())
```

### Kafka Producer

```go
import "github.com/banking/shared/kafka"

cfg := kafka.DefaultProducerConfig([]string{"localhost:9092"}, "my-service")
producer, err := kafka.NewProducer(cfg, logger)

err = producer.Publish(ctx, topic, event)
```

### Models

```go
import "github.com/banking/shared/models"

if models.StatusPending.IsFinal() {
    // ...
}

if models.IsValidCurrency("USD") {
    // Valid currency
}
```

### Validators

```go
import "github.com/banking/shared/validators"

if err := validators.ValidateTransferAmount(amount); err != nil {
    return err
}
```

## Packages

- `events/` - Kafka event definitions and topic configuration
- `kafka/` - Kafka producer and consumer with circuit breaker
- `models/` - Shared domain models (Transaction, User, Account)
- `validators/` - Input validation utilities
