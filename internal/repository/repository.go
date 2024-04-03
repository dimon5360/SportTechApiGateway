package repository

const InvalidInputParameter string = "invalid input parameter"

type AuthInterface interface {
	Init() error
	Login(interface{}) (interface{}, error)
	Register(interface{}) (interface{}, error)
	Refresh(interface{}) (interface{}, error)
}

type ProfileInterface interface {
	Init() error
	Create(interface{}) (interface{}, error)
	Read(interface{}) (interface{}, error)
	Update(interface{}) (interface{}, error)
	Delete(interface{}) error
}

type ReportInterface interface {
	Init() error
	Create(interface{}) (interface{}, error)
	Read(interface{}) (interface{}, error)
	Update(interface{}) (interface{}, error)
	Delete(interface{}) error
}
