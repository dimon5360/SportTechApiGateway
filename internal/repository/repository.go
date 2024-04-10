package repository

import "app/main/internal/dto"

const (
	InvalidInputParameter string = "invalid input parameter"
)

type AuthInterface interface {
	Init() error
	Login(req *dto.RestLoginRequest) (*dto.RestLoginResponse, error)
	Register(req *dto.RestRegisterRequest) error
	Refresh(req *dto.RestRefreshTokenRequest) (*dto.RestRefreshTokenResponse, error)
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
