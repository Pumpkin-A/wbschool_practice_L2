package main

import "fmt"

type patientStorage struct {
}

type patient struct {
	ID   string
	name string
}

func NewPatientStorage() *patientStorage {
	return &patientStorage{}
}

func (ps *patientStorage) GetPatient(ID string) (patient, error) {
	return patient{ID: "100-bbb", name: "Sam"}, nil
}

type doctorStorage struct {
}

type doctor struct {
	ID             string
	specialization string
}

func NewDoctorStorage() *doctorStorage {
	return &doctorStorage{}
}

func (ds *doctorStorage) GetDoctor(ID string) (doctor, error) {
	return doctor{ID: "999-aaa", specialization: "dentist"}, nil
}

type medicalCard struct {
	patient
	diseases string
}

type medicalCardsStorage struct{}

func NewMedicalCardsStorage() *medicalCardsStorage {
	return &medicalCardsStorage{}
}

func (mcs medicalCardsStorage) GetCard(patientID string) (medicalCard, error) {
	return medicalCard{patient: patient{ID: "100-bbb", name: "Sam"}, diseases: "some diseases"}, nil
}

type timeTable struct{}

func NewTimeTable() *timeTable {
	return &timeTable{}
}

type timeTableRecord struct {
	patient patient
	doctor  doctor
	time    string
}

func (tt *timeTable) CreateTimeTableRecord(patient patient, doctor doctor, time string) (timeTableRecord, error) {
	return timeTableRecord{patient: patient, doctor: doctor, time: time}, nil
}

type notification struct {
}

func NewNotification() *notification {
	return &notification{}
}

func (n *notification) Send(patient patient, timeTableRecord timeTableRecord) {
	fmt.Printf("%s, вы успешно записаны на прием к врачу! Дата записи: %s\n", patient.name, timeTableRecord.time)
}
