package repository

import (
	"app/main/internal/repository/reportRepository"
	"app/main/internal/repository/userRepository"
)

func NewUserRepository() userRepository.Interface {
	return userRepository.New()
}

func NewReportRepository() reportRepository.Interface {
	return reportRepository.New()
}
