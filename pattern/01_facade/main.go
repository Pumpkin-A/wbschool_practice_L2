package main

import "fmt"

// В данном примере реализован пример использования паттерна "фасад" на примере предоставления простого интерфейса
// внутренней логики сервера для сетевого слоя приложения.

// В качестве фасада выделен интерфейс HospitalTimetableFacade, который реализует структура TimeTableConcrete.
// Она содержит в себе сложную логику записи к врачу с обращением к множеству сторонних структур,
// Интерфейс представляет обрезанную(только необходимую) функциональность клиенту при обращении через http сервер
// В данном случае клиент не знает о внутренней записи пациента на прием к врачу (проверка в хранилище пациента и врача),
// поиск медицинской карты, создание записи, отправка уведомления), а просто передает необходимые данные в
// функцию записи, определенную интерфейсом, и получает данные записи.

// Плюсы: улучшение архитектуры, упрощение работы с сервисом, упрощение внесения изменений и тестирования.
// Минусы: при неправильном использовании может превратиться в божественный объект, то есть предоставлять
// разнородный, несвязный функционал.

type TimeTableConcrete struct {
	patientStorage      *patientStorage
	doctorStorage       *doctorStorage
	medicalCardsStorage *medicalCardsStorage
	notification        *notification
	timeTable           *timeTable
}

func NewTimeTableConcrete() *TimeTableConcrete {
	hospitalFacade := &TimeTableConcrete{
		patientStorage:      NewPatientStorage(),
		doctorStorage:       NewDoctorStorage(),
		medicalCardsStorage: NewMedicalCardsStorage(),
		notification:        NewNotification(),
		timeTable:           NewTimeTable(),
	}
	return hospitalFacade
}

func (tc *TimeTableConcrete) MakeAppointment(patientID string, doctorID string, timeslot string) (timeTableRecord, error) {
	//получить пациента
	patient, _ := tc.patientStorage.GetPatient(patientID)

	//получить врача
	doctor, _ := tc.doctorStorage.GetDoctor(doctorID)

	//получить медицинскую карту пациента
	tc.medicalCardsStorage.GetCard(patient.ID)

	//создать запись
	timeTableRecord, _ := tc.timeTable.CreateTimeTableRecord(patient, doctor, timeslot)

	//отправить уведомление о записи пациенту
	tc.notification.Send(patient, timeTableRecord)

	return timeTableRecord, nil
}

func (tc *TimeTableConcrete) DeleteAppointment() {}

type HTTPServer struct {
	timeTableHospital HospitalTimetableFacade
}

func NewHTTPServer(timeTableHospital HospitalTimetableFacade) *HTTPServer {
	return &HTTPServer{timeTableHospital: timeTableHospital}
}

type HospitalTimetableFacade interface {
	MakeAppointment(patientID string, doctorID string, timeslot string) (timeTableRecord, error)
	DeleteAppointment()
}

func main() {
	hospitalTimetableFacade := NewTimeTableConcrete()

	HTTPServer := NewHTTPServer(hospitalTimetableFacade)
	timeTableRecord, _ := HTTPServer.timeTableHospital.MakeAppointment("100-bbb", "999-aaa", "17.02.2025")
	fmt.Println(timeTableRecord)

}
