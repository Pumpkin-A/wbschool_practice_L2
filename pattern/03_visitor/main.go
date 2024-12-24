package main

import "fmt"

// Благодаря использованию паттерна посетитель при добавлении новых операций для иерархии классов методов оплаты
// не требуется менять сами классы.
// Плюсы: -Упрощает добавление операций, работающих со сложными структурами объектов.
// -Позволяет "не засорять" классы при необходимости добавления некоторых не связанных между собой операций,
// в том числе если они имеют смысл только для некоторых классов из существующей иерархии.
// Минусы: - Может привести к нарушению инкапсуляции элементов.
// - Не подходит, если иерархия классов часто меняется.

type PaymentMethod interface {
	getType() string
	accept(Visitor)
}

type Cash struct {
}

func (c *Cash) accept(v Visitor) {
	v.visitForCash(c)
}

func (c *Cash) getType() string {
	return "Оплата наличными"
}

type Cashless struct {
	accountNumber int
}

func (c *Cashless) accept(v Visitor) {
	v.visitForCashless(c)
}

func (c *Cashless) getType() string {
	return "Безналичный расчет"
}

type Bonuses struct {
	clientID string
}

// единственное изменение, метод, который лишь однажды нужно реализовать в каждом классе, после чего в него будет передан
// интерфейс посетителя, что позволяет добавлять любое количество новых операций, не изменяя код классов.
func (b *Bonuses) accept(v Visitor) {
	v.visitForBonuses(b)
}

func (b *Bonuses) getType() string {
	return "Оплата бонусами"
}

type Visitor interface {
	visitForCash(*Cash)
	visitForCashless(*Cashless)
	visitForBonuses(*Bonuses)
}

// Подтверждение оплаты требуется только некоторым видам оплаты. Например, при наличном расчете эта операция не нужна,
// а при оплате бонусами она будет отличаться от классической(как при проведении оплаты через банк). Таким образом,
// реализация интерфейса visitor может быть различной, в том числе и пустой.
type ConfirmationPayment struct {
}

func (cp *ConfirmationPayment) visitForCash(c *Cash) {
}

func (cp *ConfirmationPayment) visitForCashless(c *Cashless) {
	fmt.Println("Проверка подтверждения оплаты от банка")
}
func (cp *ConfirmationPayment) visitForBonuses(b *Bonuses) {
	fmt.Println("Проверка бонусов на счете в личном кабинете приложения")
}

type CreatingPayment struct {
}

func (cp *CreatingPayment) visitForCash(c *Cash) {
	fmt.Println("создание оплаты для наличного расчета")
}

func (cp *CreatingPayment) visitForCashless(c *Cashless) {
	fmt.Println("Создание оплаты для безналичного расчета")
}
func (cp *CreatingPayment) visitForBonuses(b *Bonuses) {
	fmt.Println("Создание оплаты для расчета бонусами")
}

func main() {
	cash := &Cash{}
	fmt.Println(cash.getType())
	cashless := &Cashless{}
	fmt.Println(cashless.getType())
	bonuses := Bonuses{}
	fmt.Println(bonuses.getType())
	fmt.Println()

	confirmationPayment := &ConfirmationPayment{}
	cash.accept(confirmationPayment)
	cashless.accept(confirmationPayment)
	bonuses.accept(confirmationPayment)
	fmt.Println()

	creatingPayment := &CreatingPayment{}
	cash.accept(creatingPayment)
	cashless.accept(creatingPayment)
	bonuses.accept(creatingPayment)

}
