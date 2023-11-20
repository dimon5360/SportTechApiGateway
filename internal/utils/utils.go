package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type envHandler struct {
	dict map[string]string
}

func Init() envHandler {

	var handler envHandler
	handler.dict = make(map[string]string)
	return handler
}

func (h *envHandler) Load(env_path string) {

	file, err := os.Open(env_path)
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
				panic("invalid line in env file")
			}

			h.dict[res[0]] = res[1]
		}
	}
}

func (h *envHandler) Value(key string) string {
	if val, ok := h.dict[key]; ok {
		return val
	}
	return "UNDEFINED"
}
