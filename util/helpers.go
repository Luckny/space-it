package util

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func ExtractAuthHeader(header string) (string, string, error) {
	// decode base 64 auth header
	info, err := base64.StdEncoding.DecodeString(strings.Split(header, " ")[1])
	if err != nil {
		return "", "", err
	}
	creds := strings.Split(string(info), ":")
	if len(creds) != 2 {
		return "", "", fmt.Errorf("invalid authorization header")
	}

	if len(creds[0]) == 0 || len(creds[1]) == 0 {
		return "", "", fmt.Errorf("invalid credentials")
	}

	return creds[0], creds[1], nil
}
