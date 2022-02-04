package storage

import (
	"log"
	"os"
	"path/filepath"
)

const FILENAME = "./secret.json"

func GetSecretFilePath() string {
	dirPath := getDirPath()
	return filepath.Join(dirPath, FILENAME)
}

func GetSecretFile() (*os.File, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return os.OpenFile(GetSecretFilePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	}
	return file, nil
}

func GetTokenFilePath() string {
	return filepath.Join(getDirPath(), "./token.json")
}

func GetTokenFile() (*os.File, error) {
	file, err := os.OpenFile(GetTokenFilePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		os.MkdirAll(getDirPath(), os.ModePerm)
		return os.Open("./token.json")
	}
	log.Println("HERE!")
	return file, nil
}
