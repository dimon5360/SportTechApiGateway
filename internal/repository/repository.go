package repository

type Interface interface {
	Init() error
	Add(interface{}) (interface{}, error)
	Get(interface{}) (interface{}, error)
	IsExist(interface{}) (bool, error)
	Verify(interface{}) (interface{}, error)
}
