package main

// тариф экономный
type EconomicTravelPackage struct {
	Country           string
	PlaneClassService string
	HotelClass        string
	NumExcursion      string
	Nutrition         string
}

func NewEconomicTravelPackage(country string) *EconomicTravelPackage {
	return &EconomicTravelPackage{Country: country}
}

func (ep *EconomicTravelPackage) setPlaneClassService() {
	ep.PlaneClassService = "economic"
}

func (ep *EconomicTravelPackage) setHotelClass() {
	ep.HotelClass = "hostel"
}

func (ep *EconomicTravelPackage) setNumExcursion() {
	ep.NumExcursion = "without excursion"
}

func (ep *EconomicTravelPackage) setNutrition() {
	ep.Nutrition = "without nutrition"
}

func (ep *EconomicTravelPackage) getTravelPackage() TravelPackage {
	return TravelPackage{
		Country:           ep.Country,
		PlaneClassService: ep.PlaneClassService,
		HotelClass:        ep.HotelClass,
		NumExcursion:      ep.NumExcursion,
		Nutrition:         ep.Nutrition,
	}
}

// тариф стандарт
type StandardTravelPackage struct {
	Country           string
	PlaneClassService string
	HotelClass        string
	NumExcursion      string
	Nutrition         string
}

func NewStandardTravelPackage(country string) *StandardTravelPackage {
	return &StandardTravelPackage{Country: country}
}

func (sp *StandardTravelPackage) setPlaneClassService() {
	sp.PlaneClassService = "comfort"
}

func (sp *StandardTravelPackage) setHotelClass() {
	sp.HotelClass = "three-star hotel"
}

func (sp *StandardTravelPackage) setNumExcursion() {
	sp.NumExcursion = "2"
}

func (sp *StandardTravelPackage) setNutrition() {
	sp.Nutrition = "breakfast included"
}

func (sp *StandardTravelPackage) getTravelPackage() TravelPackage {
	return TravelPackage{
		Country:           sp.Country,
		PlaneClassService: sp.PlaneClassService,
		HotelClass:        sp.HotelClass,
		NumExcursion:      sp.NumExcursion,
		Nutrition:         sp.Nutrition,
	}
}

// тариф бизнес
type BusinessTravelPackage struct {
	Country           string
	PlaneClassService string
	HotelClass        string
	NumExcursion      string
	Nutrition         string
}

func NewBusinessTravelPackage(country string) *BusinessTravelPackage {
	return &BusinessTravelPackage{Country: country}
}

func (sp *BusinessTravelPackage) setPlaneClassService() {
	sp.PlaneClassService = "business"
}

func (sp *BusinessTravelPackage) setHotelClass() {
	sp.HotelClass = "five-star hotel"
}

func (sp *BusinessTravelPackage) setNumExcursion() {
	sp.NumExcursion = "5"
}

func (sp *BusinessTravelPackage) setNutrition() {
	sp.Nutrition = "full included"
}

func (sp *BusinessTravelPackage) getTravelPackage() TravelPackage {
	return TravelPackage{
		Country:           sp.Country,
		PlaneClassService: sp.PlaneClassService,
		HotelClass:        sp.HotelClass,
		NumExcursion:      sp.NumExcursion,
		Nutrition:         sp.Nutrition,
	}
}
