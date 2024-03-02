package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type envHandler struct {
	dict          map[string]string
	isInitialized bool
}

var handler envHandler

func Env() *envHandler {

	if !handler.isInitialized {
		handler.dict = make(map[string]string)
		handler.isInitialized = true
	}
	return &handler
}

func (h *envHandler) Load(paths ...string) {

	if !h.isInitialized {
		log.Fatal("Env isn't initalized")
	}

	for _, path := range paths {

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) != 0 && !strings.HasPrefix(line, "#") {
				res := strings.Split(line, "=")

				if len(res) != 2 {
					log.Fatal("invalid line in env file")
				}

				h.dict[res[0]] = res[1]
			}
		}
	}
}

func (h *envHandler) Value(key string) string {
	if val, ok := h.dict[key]; ok {
		return val
	}
	return "UNDEFINED"
}
