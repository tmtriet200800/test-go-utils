package pkgEventBus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	pkgDomain "github.com/tmtriet200800/test-go-utils/domain"
	pkgError "github.com/tmtriet200800/test-go-utils/errors"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type EventBusKafka struct {
	logger     *logrus.Logger
	kafkaWriter *kafka.Writer
	kafkaReader *kafka.Reader
}

func NewKafka(kafkaUrl string) EventBusInterface {
	return &EventBusKafka{
		logger: logrus.New(),
		kafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(kafkaUrl),
			Topic:    "post_management",
			Balancer: &kafka.LeastBytes{},
		},
		kafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: strings.Split(kafkaUrl, ","),
			GroupID: "post_management_subscribe_group",
			Topic: "post_management",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (bus *EventBusKafka) Publish(parentCtx context.Context, event *pkgDomain.Event) error {
	eventPayload, err := json.Marshal(event)

	if err != nil {
		return pkgError.Wrap(err)
	}
	
	msg := kafka.Message{
		Key: event.ID[:],
		Value: eventPayload,
	}	

	if err = bus.kafkaWriter.WriteMessages(parentCtx, msg); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func (bus *EventBusKafka) Subscribe(ctx context.Context, eventType string, fn EventHandler) error {
	for {
		msg, errRead := bus.kafkaReader.ReadMessage(ctx)

		bus.logger.Info(fmt.Sprintf("[EventBus] Receiving: %v", msg))

		
		if errRead != nil {
			return pkgError.Wrap(errRead)
		}

		var domainEvent pkgDomain.Event

		if errParse := json.Unmarshal([]byte(msg.Value), &domainEvent); errParse != nil {
			return pkgError.Wrap(errParse)
		}

		// bus.logger.Info(fmt.Sprintf("[EventBus] Receiving: %v", domainEvent))

		if domainEvent.Type == eventType {
			return fn(ctx, &domainEvent)
		}
	}
}

