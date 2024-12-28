package main

import (
	"errors"
	"fmt"
)

// "Цепочка обязанностей" реализована на примере операции создания оплаты заказа.
//  Паттерн позволяет переложить обязанность отслеживания выполненных условий на конкретных обработчиков,
// реализуя принцип единой ответственности.
// Все конкретные обработчики реализуют интерфейс OrderDetail, это позволяет строить нужную цепочку динамически

// Плюсы: - реализован принцип единой ответственности
// - реализован принцип открытости/закрытости
// - уменьшает зависимость между клиентом и обработчиками.
// Минусы: Запрос может остаться никем не обработанным.

type OrderDetail interface {
	execute(*CreatingOrder) error
	setNext(OrderDetail)
}

// начальный элемент цепочки
type NewOrderPricingChain struct {
	next OrderDetail
}

func (op *NewOrderPricingChain) execute(co *CreatingOrder) error {
	op.setNext(&Authorization{})
	return op.next.execute(co)
}

func (r *NewOrderPricingChain) setNext(next OrderDetail) {
	r.next = next
}

type Authorization struct {
	next OrderDetail
}

func (r *Authorization) execute(co *CreatingOrder) error {
	if co.clientUid != "" {
		fmt.Println("клиент авторизован")
		r.setNext(&CheckPromoCode{})
		return r.next.execute(co)
	}
	return errors.New("client not authorized")
}

func (r *Authorization) setNext(next OrderDetail) {
	r.next = next
}

type CheckPromoCode struct {
	next OrderDetail
}

// Реализовано динамическое построение цепочки: в зависимости от того, указан ли промокод, выбирается необходимый обработчик.
// хотя динамическое построение необязательно, порядок обработчиков можно задать и статически, в клиенте
func (cpc *CheckPromoCode) execute(co *CreatingOrder) error {
	if co.PromoCode != "" {
		fmt.Println("указан промокод")
		cpc.setNext(&ApplyPromoCode{})
		return cpc.next.execute(co)
	}
	fmt.Println("промокод не указан")
	cpc.setNext(&CreatePayment{})
	return cpc.next.execute(co)
}

func (d *CheckPromoCode) setNext(next OrderDetail) {
	d.next = next
}

type ApplyPromoCode struct {
	next OrderDetail
}

func (apc *ApplyPromoCode) execute(co *CreatingOrder) error {
	co.TotalPrice -= 10
	fmt.Println("промокод активирован")
	apc.setNext(&CreatePayment{})
	return apc.next.execute(co)
}

func (d *ApplyPromoCode) setNext(next OrderDetail) {
	d.next = next
}

type CreatePayment struct {
	next OrderDetail
}

func (pm *CreatePayment) execute(co *CreatingOrder) error {
	if co.PaymentMethod != "" {
		fmt.Println("метод оплаты выбран")
		return nil
	}
	return errors.New("payment method not set")
}

func (pm *CreatePayment) setNext(next OrderDetail) {
	pm.next = next
}

type CreatingOrder struct {
	clientUid     string
	OrderUid      string
	TotalPrice    int
	PaymentMethod string
	PromoCode     string
}

func main() {
	orderPricing := &NewOrderPricingChain{}

	creatingOrder := &CreatingOrder{
		clientUid:     "123-rtt",
		OrderUid:      "567-fgh",
		TotalPrice:    100,
		PaymentMethod: "wb-pay",
		PromoCode:     "NEW_YEAR",
	}

	err := orderPricing.execute(creatingOrder)
	if err != nil {
		fmt.Println(err)
	}
}
