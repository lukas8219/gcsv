package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const tempFile = "temp-config.txt"
const configDirName = "/.gcsv"
const configFilePath = "config.txt"
const separator = "="

type Storage struct {
	Conn *sql.Conn
}

var storage *Storage

type FavoriteSheet struct {
	Name string
	ID   string
}

func GetStorage() Storage {
	if storage == nil {
		storage = &Storage{
			Conn: GetConnection(),
		}
	}
	return *storage
}

func getDirPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Not able to get user home dir", err)
	}
	return filepath.Join(dirname, configDirName)
}

func (s *Storage) Store(sheet FavoriteSheet) {
	query := `
	INSERT INTO favorites(name, sheetId) VALUES('%s','%s');
	`
	fmtInsert := fmt.Sprintf(query, sheet.Name, sheet.ID)
	_, err := s.Conn.ExecContext(context.Background(), fmtInsert)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Saved %s with ID: %s", sheet.Name, sheet.ID))
}

func (s *Storage) Remove(id string) error {
	fmtQuery := fmt.Sprintf(`
	DELETE FROM favorites WHERE name = '%s';
	`, id)
	_, err := s.Conn.ExecContext(context.Background(), fmtQuery)
	return err
}

func (s *Storage) ListAll() []FavoriteSheet {
	list := make([]FavoriteSheet, 0)
	query := `
	SELECT name, sheetId FROM favorites;
	`
	res, err := s.Conn.QueryContext(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() == true {
		var ref FavoriteSheet
		res.Scan(&ref.Name, &ref.ID)
		list = append(list, ref)
	}

	return list
}

func (s *Storage) Get(name string) (string, error) {
	query := fmt.Sprintf(`
	SELECT sheetId FROM favorites WHERE name = '%s';
	`, name)

	res, err := s.Conn.QueryContext(context.Background(), query)
	if err != nil {
		return "", err
	}

	if res.Next() == true {
		var id string
		res.Scan(&id)
		return id, nil
	} else {
		return "", errors.New("Not found")
	}
}
