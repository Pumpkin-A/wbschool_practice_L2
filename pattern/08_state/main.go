package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron"
)

// Паттерн состояние в данном примере реализован для сущности Order, состояние которого меняется по мере
// прогресса работы над заказом. Изменение состояний происходит внутри них самих. Реализация методов для каждого состояния
// отображает поведение заказа в каждом из этих состояний.

// Плюсы: Избавляет от множества больших условных операторов машины состояний.
//  Концентрирует в одном месте код, связанный с определённым состоянием.
// Минусы: Может неоправданно усложнить код, если состояний мало и они редко меняются.
// Если валидных переходов между состояниями мало, то необходимо создавать методы состояний, которые просто возвращают ошибку

type Order struct {
	currentState State
}

func NewOrder() *Order {
	order := &Order{}
	order.currentState = &StateCreated{order: order}
	return order
}

// начало бронирования возможно только из состояния stateCreated, это регулируется внутри самого состояния.
// Аналогично для оплаты (только из состояния stateBooked)
func (o *Order) CreateBooking() error {
	fmt.Println("клиент нажал кнопку забронировать заказ")
	return o.currentState.Book()
}

func (o *Order) CreatePayment() error {
	fmt.Println("клиент нажал кнопку оплатить заказ")
	return o.currentState.Pay()
}

func (o *Order) SetState(state State) {
	o.currentState = state
}

type State interface {
	SetState(state State)
	Book() error
	Pay() error
}

type StateCreated struct {
	order *Order
}

func (c *StateCreated) SetState(state State) {
	c.order.currentState = state
}

// имитация процесса бронирования, по истечении которой также изменяется состояние.
func (c *StateCreated) Book() error {
	c.SetState(&StateBooking{order: c.order})
	fmt.Println("выполняется бронирование заказа")
	time.Sleep(time.Second * 3)
	fmt.Println("бронирование заказа выполнено, меняем состояние на 'забронирован'")
	c.SetState(&StateBooked{order: c.order})
	return nil
}

func (c *StateCreated) Pay() error {
	return fmt.Errorf("order out of paying state")
}

type StateBooking struct {
	order *Order
}

func (b *StateBooking) SetState(state State) {
	b.order.currentState = state
}

func (b *StateBooking) Book() error {
	return fmt.Errorf("Order is already booking")
}

func (p *StateBooking) Pay() error {
	return fmt.Errorf("Order out of paying state")
}

type StateBooked struct {
	order *Order
}

func (b *StateBooked) Book() error {
	return fmt.Errorf("order is already booked")
}

func (b *StateBooked) Pay() error {
	b.SetState(&StatePaying{order: b.order})
	fmt.Println("выполняется оплата заказа")
	time.Sleep(time.Second * 3)
	fmt.Println("оплата заказа выполнена, меняем состояние на 'оплачен'")
	b.SetState(&StatePaid{order: b.order})
	return nil
}

func (b *StateBooked) SetState(state State) {
	b.order.currentState = state
}

type StatePaying struct {
	order *Order
}

func (p *StatePaying) Book() error {
	return fmt.Errorf("order is already booked")
}

func (p *StatePaying) Pay() error {
	return fmt.Errorf("Order is already paying")
}

func (p *StatePaying) SetState(state State) {
	p.order.currentState = state
}

type StatePaid struct {
	order *Order
}

func (p *StatePaid) Book() error {
	return fmt.Errorf("order is already booked")
}

func (p *StatePaid) Pay() error {
	return fmt.Errorf("order is already paid")
}

func (p *StatePaid) SetState(state State) {
	p.order.currentState = state
}

var wg = sync.WaitGroup{}

// Имитация оплаты. Крон пытается "нажать кнопку оплаты", достигает успеха только при состоянии stateBooked.
func startCron(order *Order) {
	c := cron.New()

	c.AddFunc("@every 1s", func() {
		err := order.CreatePayment()
		fmt.Printf("Текущее состояние заказа: %T\n", order.currentState)

		if err == nil {
			wg.Done()
			c.Stop()
		}

	})
	c.Start()
}

func main() {
	wg.Add(1)
	order := NewOrder()
	startCron(order)

	// должно выводить ошибку, тк заказ еще не забронирован
	err := order.CreatePayment()
	if err != nil {
		fmt.Println(err)
	}

	// после выполнения бронирования крон отследит изменение состояния и запустит оплату самостоятельно
	err = order.CreateBooking()
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()
	fmt.Printf("Текущее состояние заказа: %T\n", order.currentState)

}
