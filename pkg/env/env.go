package env

import (
	gt "github.com/joho/godotenv"
	"log"
)

func Init() {
	err := gt.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func LoadFile(filepath ...string) {

	err := gt.Load(filepath...)
	if err != nil {
		log.Fatal(err)
	}
}
