package main

import "fmt"

// В качестве примера реализации паттерна "команда" изображен пример взаимодействия бизнес логики и базы данных.
// Для конкретных действий (создание записи в бд, запрос к бд) используются конкретные команды, это позволяет
// скрыть реализацию бд от бизнес логики, избежать копирования кода (в случае, если в программе, например,
// появится кэш, который так же будет посылать запросы в бд, когда кэш станет грязным)
// Также есть возможность сохранять историю, так как все необходимые данные для выполнения команды хранятся в ней.

// Так как в го есть "утиная типизация", то паттерн больше заключается в сохранении параметров функции в отдельной структуре

// Плюсы: Позволяет реализовать простую отмену и повтор операций.
//  Позволяет реализовать отложенный запуск операций.
// Убирает прямую зависимость между объектами, вызывающими операции, и объектами, которые их непосредственно выполняют.
// Реализует принцип открытости/закрытости.

// Минусы: Усложняет код программы из-за введения множества дополнительных классов.

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

// методы записи и чтения объявляются только для конкретной бд и используются командами
func (db *AnyDB) write(req string) error {
	fmt.Printf("Данные успешно записаны в бд: %s\n", req)
	return nil
}

func (db *AnyDB) read(req string) string {
	resp := "прочитанные данные"
	fmt.Printf("Данные прочитаны из бд. Запрос: %s, Ответ от бд: %s\n", req, resp)
	return resp
}

type BusinessLogic struct {
	db DB
}

func (b *BusinessLogic) writeRequest(recordingData string) {
	writeCommand := NewWriteCommand(b.db, recordingData)
	writeCommand.execute()
}

func (b *BusinessLogic) readRequest(req string) {
	readCommand := NewReadCommand(b.db, req)
	readCommand.execute()
}

type Cache struct {
	db DB
}

func (b *Cache) writeRequest(recordingData string) {
	writeCommand := NewWriteCommand(b.db, recordingData)
	writeCommand.execute()
}

func (b *Cache) readRequest(req string) {
	readCommand := NewReadCommand(b.db, req)
	readCommand.execute()
}

func main() {
	anyDB := &AnyDB{}

	businessLogic := &BusinessLogic{db: anyDB}

	businessLogic.writeRequest("AAA")
	businessLogic.readRequest("select * where id = 1234")
}
