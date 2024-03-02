package repository

import (
	profile "app/main/internal/repository/profile"
	report "app/main/internal/repository/report"
	user "app/main/internal/repository/user"
)

type Interface interface {
	Init() error
	Add(interface{}) (interface{}, error)
	Get(interface{}) (interface{}, error)
}

func Users() Interface {
	return user.NewUserRepository()
}

func Profiles() Interface {
	return profile.NewProfileRepository()
}

func Reports() Interface {
	return report.NewReportRepository()
}
