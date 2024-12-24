package main

import "fmt"

// В качестве примера реализации паттерна "команда" изображен пример взаимодействия бизнес логики и базы данных.
// Для конкретных действий (создание записи в бд, запрос к бд) используются конкретные команды, это позволяет
// скрыть реализацию бд от бизнес логики, избежать копирования кода (в случае, если в программе, например,
// появится кэш, который так же будет посылать запросы в бд, когда кэш станет грязным)

type BusinessLogic struct {
	command Command
}

func (b *BusinessLogic) request() {
	b.command.execute()
}

type Command interface {
	execute()
}

type WriteCommand struct {
	db            DB
	recordingData string
}

func NewWriteCommand(db DB, recordingData string) *WriteCommand {
	return &WriteCommand{
		db:            db,
		recordingData: recordingData,
	}
}

func (c *WriteCommand) execute() {
	c.db.write(c.recordingData)
}

type ReadCommand struct {
	db   DB
	req  string
	resp string
}

func NewReadCommand(db DB, req string) *ReadCommand {
	return &ReadCommand{
		db:  db,
		req: req,
	}
}

func (c *ReadCommand) execute() {
	resp := c.db.read(c.req)
	c.resp = resp
}

type DB interface {
	write(req string) error
	read(req string) string
}

type AnyDB struct {
}

func (db *AnyDB) write(req string) error {
	fmt.Printf("Данные успешно записаны в бд: %s\n", req)
	return nil
}

func (db *AnyDB) read(req string) string {
	resp := "прочитанные данные"
	fmt.Printf("Данные прочитаны из бд. Запрос: %s, Ответ от бд: %s\n", req, resp)
	return resp
}

func main() {
	anyDB := &AnyDB{}

	writeCommand := NewWriteCommand(anyDB, "newData")
	businessLogic := &BusinessLogic{command: writeCommand}
	businessLogic.command.execute()

	readCommand := NewReadCommand(anyDB, "getAnyData")
	businessLogic.command = readCommand
	businessLogic.command.execute()
	fmt.Printf("resp: %s", readCommand.resp)
}
