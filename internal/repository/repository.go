package repository

import "app/main/internal/dto"

const InvalidInputParameter string = "invalid input parameter"

type AuthInterface interface {
	Init() error
	Login(req *dto.RestAuthRequest) (*dto.RestLoginResponse, error)
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
