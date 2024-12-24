package mymain

type DataManipulator interface {
	openConnection() ([]Event, error)
	closeConnction([]Event) error
}
