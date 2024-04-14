package repository

import (
	"app/main/internal/repository/authRepository"
	"app/main/internal/repository/profileRepository"
	"app/main/internal/repository/reportRepository"
)

func NewAuthRepository() authRepository.Interface {
	return authRepository.New()
}

func NewProfileRepository() profileRepository.Interface {
	return profileRepository.New()
}

func NewReportRepository() reportRepository.Interface {
	return reportRepository.New()
}
