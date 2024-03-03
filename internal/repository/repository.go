package repository

type Interface interface {
	Init() error
	Add(interface{}) (interface{}, error)
	Get(interface{}) (interface{}, error)
}
