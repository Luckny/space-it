package util

import (
	"log"
	"os"

	"github.com/google/uuid"
)

var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func GenUUID() uuid.UUID {
	return uuid.New()
}
