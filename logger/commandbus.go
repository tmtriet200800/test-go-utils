package pkgLogger

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type commandbusLogger struct {
	logger *logrus.Logger
	db *sqlx.DB
	tablename string
}

type commandBusLog struct {
	owner string
	action string
	command string
	payload string
}

// type commandBusLogDB struct {
// 	owner string `db:"owner"`
// 	action string `db:"action"`
// 	command string `db:"command"`
// 	level string `db:"level"`
// 	payload string `db:"payload"`
// }

func NewCommandBusLogger(db *sqlx.DB, tablename string) LoggerInterface{
	return &commandbusLogger{
		logger: logrus.New(),
		db: db,
		tablename: tablename,
	}	
}

func NewCommandBusLog(owner string, action string, command string, payload string) interface{}{
	return commandBusLog{
		owner: owner,
		action: action,
		command: command,
		payload: payload,
	}
}

func (l *commandbusLogger) SaveToDB(commandBusLog commandBusLog, logType string){
	l.db.QueryRowx(
		fmt.Sprintf(
			"INSERT INTO %v (owner, action, command, level, payload) VALUES ($1, $2, $3, $4, $5) RETURNING *",
			l.tablename, 
		),
		commandBusLog.owner,
		commandBusLog.action,
		commandBusLog.command,
		logType,		
		commandBusLog.payload,
	)
}

func (l *commandbusLogger) Info(verbose bool, data interface{}){
	commandBusLog := data.(commandBusLog)

	if verbose{
		logrus.Info(fmt.Sprintf("[CommandBus - %v - %v] %v", commandBusLog.owner, commandBusLog.command, commandBusLog.action))
	}

	l.SaveToDB(commandBusLog, "INFO")
}

func (l *commandbusLogger) Debug(verbose bool, data interface{}){
	commandBusLog := data.(commandBusLog)

	if verbose{
		logrus.Debug(fmt.Sprintf("[CommandBus - %v - %v] %v", commandBusLog.owner, commandBusLog.command, commandBusLog.action))
	}

	l.SaveToDB(commandBusLog, "DEBUG")
}

func (l *commandbusLogger) Error(verbose bool, data interface{}){
	commandBusLog := data.(commandBusLog)

	if verbose{
		logrus.Error(fmt.Sprintf("[CommandBus - %v - %v] %v", commandBusLog.owner, commandBusLog.command, commandBusLog.action))
	}

	l.SaveToDB(commandBusLog, "ERROR")
}
