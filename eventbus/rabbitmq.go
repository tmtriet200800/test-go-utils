package pkgEventBus

import (
	"context"
	"encoding/json"
	"fmt"
	pkgDomain "go_utils/domain"
	pkgError "go_utils/errors"
	"math/rand"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type EventBusRabbitMQ struct {
	logger     *logrus.Logger
	channel *amqp.Channel
	mainQueue amqp.Queue
	callBackQueue amqp.Queue
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
			bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}


func NewRabbitMQ() EventBusInterface{
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	
	if err != nil{
		logrus.New().Error("Failed to connect to RabbitMQ")
		panic(err)
	}

	ch, err := conn.Channel()

	if err != nil{
		logrus.New().Error("Failed to open a channel")
		panic(err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil{
		logrus.New().Error("Failed to set QoS")
		panic(err)
	}

	mainQueue, err := ch.QueueDeclare(
		"event-bus", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil{
		logrus.New().Error("Failed to declare a queue")
		panic(err)
	}

	callBackQueue, err := ch.QueueDeclare(
		"event-bus-callback", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil{
		logrus.New().Error("Failed to declare a callback queue")
		panic(err)
	}

	return &EventBusRabbitMQ{
		logger: logrus.New(),
		mainQueue: mainQueue,
		callBackQueue: callBackQueue,
		channel: ch,
	}
}


func (bus *EventBusRabbitMQ) Publish(parentCtx context.Context, event *pkgDomain.Event) error {	
	eventPayload, errParse := json.Marshal(event)

	if errParse != nil {return pkgError.Wrap(errParse)}

	msgsCallBack, err := bus.channel.Consume(
		bus.callBackQueue.Name, // queue
		"callBackConsumer",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil{
		return pkgError.New("[EventBus - Publisher] Failed to register a callback queue")
	}

	corrId := randomString(32)

	if err:= bus.channel.Publish(
		"",
		bus.mainQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			CorrelationId: corrId,
			Body: eventPayload,
			ReplyTo: bus.callBackQueue.Name,
		},
	); err != nil {
		return pkgError.Wrap(err)
	}

	var errLogicString string

	bus.logger.Info("[EventBus - Publisher] Published a message: ", event.Payload)

	for d := range msgsCallBack {
		if corrId == d.CorrelationId {
			bus.logger.Info("[EventBus - Publisher] Recevied response from subscriber")
			errLogicString = string(d.Body[:])
			break
		}

		bus.logger.Info("[EventBus - Publisher] Ignore mismatch response from subscriber")
	}

	bus.channel.Cancel("callBackConsumer", false)

	if errLogicString != "" {
		return pkgError.New(errLogicString)
	}

	return nil
}



func (bus *EventBusRabbitMQ) Subscribe(ctx context.Context, eventType string, fn EventHandler) error{
	msgs, err := bus.channel.Consume(
		bus.mainQueue.Name,
		"",
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil{
		return pkgError.New("[EventBus - Subscriber] Failed to register a consumer")
	}

	go func() {
		for d := range msgs {
			var domainEvent pkgDomain.Event
	
			if err := json.Unmarshal(d.Body, &domainEvent); err != nil{
				bus.logger.Error(fmt.Sprintf("[EventBus - Subscriber] Error in parsing: %v", err))
			}

			bus.logger.Info(fmt.Sprintf("[EventBus - Subcriber] Received a event: %v with payload %v", domainEvent.Type, domainEvent.Payload))

			if domainEvent.Type == eventType {
				errLogicString := ""

				if errLogic := fn(ctx, &domainEvent); errLogic != nil {
					bus.logger.Error("[EventBus - Subscriber] ", errLogic)
					errLogicString = errLogic.Error()
				}

				err = bus.channel.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          []byte(errLogicString),
				})

				bus.logger.Info("[EventBus - Subcriber] Published response message")
			
			}		
		}
	}()

	return nil
}
