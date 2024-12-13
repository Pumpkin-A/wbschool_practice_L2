package main

import "fmt"

// Выполнена реализация паттерна "строитель" с помощью интерфейса TravelPackageBuilder,
// структур с разными типами класса обслуживания, который его реализуют и директора, отвечающим за компоновку поездки.

// В данном случае "директор" был добавлен не потому, что порядок компоновки пакетов обслуживания отличается,
// а для наглядной демонстрации использования интерфейса и сокрытия от клиентского кода процесс конструирования объектов.

// Так как страна - место назначения путешествия может быть любым, она явно передается в конструктор каждого пакета,
// для предотвращения создания новых билдеров для каждой новой страны.

// Плюсы: изоляция сложный кода сборки продукта от его основной бизнес-логики, пошаговое создание продуктов.
// Минусы: - Усложняет код программы из-за введения дополнительных классов.
// 			- Клиент будет привязан к конкретным классам строителей, так как в интерфейсе директора
// может не быть метода получения результата.
// -при необходимости увеличения вариативности создаваемых объектов паттерн становится менее полезным.
// (Например, если в компании по продаже путевок захотят провести акцию - при поезде в Египет плюс две бесплатные экскурсии
// нужно заносить данные в каждого "билдера")

type TravelPackage struct {
	Country           string
	PlaneClassService string
	HotelClass        string
	NumExcursion      string
	Nutrition         string
}

type TravelPackageBuilder interface {
	setPlaneClassService()
	setHotelClass()
	setNumExcursion()
	setNutrition()
	getTravelPackage() TravelPackage
}

func getTravelPackageBuilder(travelPackageType, country string) TravelPackageBuilder {
	if travelPackageType == "economic" {
		return NewEconomicTravelPackage(country)
	}

	if travelPackageType == "standard" {
		return NewStandardTravelPackage(country)
	}

	if travelPackageType == "business" {
		return NewBusinessTravelPackage(country)
	}
	return nil
}

type Director struct {
	TravelPackageBuilder TravelPackageBuilder
}

func newDirector(tpBuilder TravelPackageBuilder) *Director {
	return &Director{
		TravelPackageBuilder: tpBuilder,
	}
}

func (d *Director) SetTravelPackageBuilder(tpBuilder TravelPackageBuilder) {
	d.TravelPackageBuilder = tpBuilder
}

func (d *Director) OrganizeTrip() TravelPackage {
	d.TravelPackageBuilder.setPlaneClassService()
	d.TravelPackageBuilder.setHotelClass()
	d.TravelPackageBuilder.setNumExcursion()
	d.TravelPackageBuilder.setNutrition()
	return d.TravelPackageBuilder.getTravelPackage()
}

func main() {
	businessTravelPackage := getTravelPackageBuilder("business", "Belarus")
	economicTravelPackage := getTravelPackageBuilder("economic", "Belarus")

	director := newDirector(businessTravelPackage)
	businessTrip := director.OrganizeTrip()

	fmt.Printf("Business trip plane class service: %s\n", businessTrip.PlaneClassService)
	fmt.Printf("Business trip hotel class: %s\n", businessTrip.HotelClass)
	fmt.Printf("Business trip plane nutrition: %s\n", businessTrip.Nutrition)
	fmt.Printf("Business trip plane number of excursions: %s\n\n", businessTrip.NumExcursion)

	director.SetTravelPackageBuilder(economicTravelPackage)
	economicTrip := director.OrganizeTrip()

	fmt.Printf("Economic trip plane class service: %s\n", economicTrip.PlaneClassService)
	fmt.Printf("Economic trip hotel class: %s\n", economicTrip.HotelClass)
	fmt.Printf("Economic trip plane nutrition: %s\n", economicTrip.Nutrition)
	fmt.Printf("Economic trip plane number of excursions: %s\n", economicTrip.NumExcursion)

}
