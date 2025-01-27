package main

import (
	"errors"
	"fmt"
)

// Паттерн "фабричный метод" продемонстрирован на примере реализации разными провайдерами интерфейса TransactionProviderConcrete.
// Использование паттерна полезно в данном примере, так как провайдеры транзакций потенциально могут быть использованы в
// различных частях программы, их количество может постоянно увеличиваться. Также клиент не знает о внутренней реализации
// методов и различиях в работе провайдеров, ему необходимо просто выбрать нужный (принцип открытости/закрытости)

// Плюсы: избавление от привязки к конкретным классам; конструктор и вся реализация в одном месте, что упрощает поддержку кода;
// упрощает добавление новых продуктов в программу; реализует принцип открытости/закрытости
// Минусы: Если использовать упрощенную версию фабричного метода, мы избавляемся от его главной проблемы
// "Может привести к созданию больших параллельных иерархий классов, так как для каждого класса
// продукта надо создать свой подкласс создателя."

// Проведение транзакций через сбп
type SBPTransactionProvider struct {
}

func NewSBP() *SBPTransactionProvider {
	return &SBPTransactionProvider{}
}

func (sbp *SBPTransactionProvider) makePayment(transactionID string) {
	fmt.Printf("списание через СБП по транзакции: %s\n", transactionID)
}

func (sbp *SBPTransactionProvider) makeRefund(transactionID string) {
	fmt.Printf("возврат через СБП по транзакции: %s\n", transactionID)
}

// Проведение транзакций через wbpay
type WBpayTransactionProvider struct {
}

func NewWBpay() *WBpayTransactionProvider {
	return &WBpayTransactionProvider{}
}

func (wbpay *WBpayTransactionProvider) makePayment(transactionID string) {
	fmt.Printf("списание через wbpay по транзакции: %s\n", transactionID)
}

func (wbpay *WBpayTransactionProvider) makeRefund(transactionID string) {
	fmt.Printf("возврат через wbpay по транзакции: %s\n", transactionID)
}

// объявление интерфейса провайдера
type TransactionProviderConcrete interface {
	makePayment(transactionID string)
	makeRefund(transactionID string)
}

// Это "простой" фабричный метод, одна из двух вариаций паттерна, где класс creator - конкретный.
type Creator struct {
}

func (c *Creator) createProvider(providerName string) (TransactionProviderConcrete, error) {
	switch providerName {
	case "sbp":
		return NewSBP(), nil
	case "wbpay":
		return NewWBpay(), nil
	default:
		return nil, errors.New("unknown type of provider")
	}
}

func main() {
	creator := Creator{}
	sbpProvider, _ := creator.createProvider("sbp")
	sbpProvider.makePayment("1234-asdf")
	sbpProvider.makeRefund("1234-asdf")
	fmt.Println()
	wbpayProvider, _ := creator.createProvider("wbpay")
	wbpayProvider.makePayment("456-sdf")
	wbpayProvider.makeRefund("4567-sdfg")
}
