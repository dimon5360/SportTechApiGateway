package env

import (
	gt "github.com/joho/godotenv"
)

func Init() error {
	return gt.Load()
}
