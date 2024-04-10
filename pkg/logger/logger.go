package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

func Init() error {

	filename := fmt.Sprintf("./%s/%s.log", os.Getenv("LOGS_FOLDER"),
		time.Now().Format(os.Getenv("LOGS_FORMAT_DATETIME")))

	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println(sig.String())
			logFile.Close() // close output file after application termination
			os.Exit(0)
		}
	}()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Logger initialized")

	return nil
}
