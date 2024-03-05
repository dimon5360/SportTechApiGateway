package service

type Interface interface {
	Init() error
	Run() error
}
