package pkgCommandBus

import (
	"context"
	"encoding/json"
	"math/rand"

	pkgDomain "github.com/tmtriet200800/test-go-utils/domain"
	pkgError "github.com/tmtriet200800/test-go-utils/errors"
	pkgLogger "github.com/tmtriet200800/test-go-utils/logger"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// CommandBusRabbitMQ allows to subscribe/dispatch commands
// Subscribing to the same command twice will unsubscribe previous handler
// command handler should be one to one
type CommandBusRabbitMQ struct {
	logger pkgLogger.LoggerInterface
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

func NewRabbitMQ(maxConcurrentCalls int, db *sqlx.DB) CommandBusInterface {
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
		"command-bus", // name
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
		"command-bus-callback", // name
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


	return &CommandBusRabbitMQ{
		logger: pkgLogger.NewCommandBusLogger(db, "log_command_bus"),
		mainQueue: mainQueue,
		callBackQueue: callBackQueue,
		channel: ch,
	}
}

func (bus *CommandBusRabbitMQ) Publish(ctx context.Context, command *pkgDomain.Command) error{
	commandPayload, errParse := json.Marshal(command)

	if errParse != nil{
		return pkgError.Wrap(errParse)
	}


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
		return pkgError.New("[CommandBus - Publisher] Failed to register a callback queue")
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
			Body: commandPayload,
			ReplyTo: bus.callBackQueue.Name,
		},
	); err != nil {
		return pkgError.Wrap(err)
	}

	var errLogicString string

	bus.logger.Info(true, pkgLogger.NewCommandBusLog("publisher", "Publish message to subscriber", command.Name, string(command.Payload)))

	for d := range msgsCallBack {
		if corrId == d.CorrelationId {
			bus.logger.Info(true, pkgLogger.NewCommandBusLog("publisher", "Recevied response from subscriber", command.Name, ""))
			errLogicString = string(d.Body[:])
			break
		}
	}

	bus.channel.Cancel("callBackConsumer", false)

	if errLogicString != "" {
		return pkgError.New(errLogicString)
	}

	return nil
}

func (bus *CommandBusRabbitMQ) Subscribe(ctx context.Context, commandName string, fn CommandHandler) error{
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
		return pkgError.New("[CommandBus - Subscriber] Failed to register a consumer")
	}

	go func() {
		for d := range msgs {
			var domainCommand pkgDomain.Command
			
			if err := json.Unmarshal(d.Body, &domainCommand); err != nil{
				bus.logger.Info(true, pkgLogger.NewCommandBusLog("subscriber", "Throw unmarshal error", domainCommand.Name, "Cannot unmarshal command payload"))
			}

			bus.logger.Info(true, pkgLogger.NewCommandBusLog("subscriber", "Recevied command from publisher", domainCommand.Name, string(domainCommand.Payload)))

			if domainCommand.Name == commandName {
				errLogicString := ""

				if errLogic := fn(ctx, &domainCommand); errLogic != nil {
					errLogicString = errLogic.Error()
					bus.logger.Info(true, pkgLogger.NewCommandBusLog("subscriber", "Throw logic error", domainCommand.Name, errLogicString))
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

				bus.logger.Info(true, pkgLogger.NewCommandBusLog("subscriber", "Publish response to publisher", domainCommand.Name, ""))
			}		
		}
	}()

	return nil
}

// func (bus *CommandBusRabbitMQ) Unsubscribe(ctx context.Context, commandName string) error {

// }


