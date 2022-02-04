package storage

import (
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
		return os.Open(GetSecretFilePath())
	}
	return file, nil
}

func GetTokenFilePath() string {
	return filepath.Join(getDirPath(), "./token.json")
}

func CreateOrWriteTokenFile() (*os.File, error) {
	return getTokenFile(os.O_RDWR | os.O_CREATE)
}

func getTokenFile(permission int) (*os.File, error) {
	file, err := os.OpenFile(GetTokenFilePath(), permission, 0600)
	if err != nil {
		os.MkdirAll(getDirPath(), os.ModePerm)
		return os.Open("./token.json")
	}
	return file, nil
}

func ReadTokenFile() (*os.File, error) {
	return getTokenFile(os.O_RDONLY)
}
