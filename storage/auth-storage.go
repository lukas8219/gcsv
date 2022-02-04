package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
)

const FILENAME = "./secret.json"

func GetSecretFilePath() string {
	dirPath := getDirPath()
	return filepath.Join(dirPath, FILENAME)
}

func GetToken() (*oauth2.Token, error) {
	con := GetConnection()
	query := `SELECT json FROM authentication WHERE type = 'token'`
	res, err := con.QueryContext(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	var result string
	res.Next()
	if err = res.Scan(&result); err != nil {
		return nil, err
	}
	var token *oauth2.Token
	json.NewDecoder(strings.NewReader(result)).Decode(&token)
	return token, err
}

func SaveToken(token string) error {
	con := GetConnection()

	statement := `UPDATE authentication
	SET json='%s'
	WHERE type='token';
	
	INSERT INTO authentication (type, json)
	SELECT 'token', '%s'
	WHERE (Select Changes() = 0);
	`

	fmtStatement := fmt.Sprintf(statement, token, token)

	_, err := con.ExecContext(context.Background(), fmtStatement)
	return err
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
	return getTokenFile(os.O_RDWR | os.O_CREATE | os.O_APPEND)
}

func getTokenFile(permission int) (*os.File, error) {
	file, err := os.OpenFile(GetTokenFilePath(), permission, 0600)
	if err != nil {
		log.Println(err)
		os.MkdirAll(getDirPath(), os.ModePerm)
		return os.Open("./token.json")
	}
	return file, nil
}

func ReadTokenFile() (*os.File, error) {
	return getTokenFile(os.O_RDONLY)
}
