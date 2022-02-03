package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const tempFile = "temp-config.txt"
const dir = "./storage"
const configFilePath = "./storage/config.txt"
const separator = "="

type FavoriteSheet struct {
	Name string
	ID   string
}

func createOrReadFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.FileMode(os.O_APPEND))
	if err != nil {
		f, err = os.Create(path)
	}
	return f, err
}

func Store(sheet FavoriteSheet) {

	_, err := Get(sheet.Name)

	if err == nil {
		log.Fatalln("There's already a Entry with this Name")
		return
	}

	f, err := createOrReadFile(configFilePath)
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s%s%s", sheet.Name, separator, sheet.ID))
	f.WriteString("\n")
	err = f.Sync()

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println(fmt.Sprintf("Saved %s with ID: %s", sheet.Name, sheet.ID))
}

func Remove(id string) error {
	temp, err := ioutil.TempFile(dir, tempFile)
	defer temp.Close()
	if err != nil {
		return err
	}

	f, err := os.Open(configFilePath)
	defer f.Close()
	if err != nil {
		return err
	}

	file := bufio.NewScanner(f)

	for file.Scan() {
		result := file.Text()
		parsed := strings.Split(result, separator)

		if len(parsed) != 2 {
			return errors.New("Unmatched config")
		}

		if parsed[0] != id {
			temp.WriteString(result)
			temp.WriteString("\n")
		}
	}
	return os.Rename(temp.Name(), configFilePath)
}

func ListAll() []FavoriteSheet {
	list := make([]FavoriteSheet, 0)

	f, err := os.Open(configFilePath)
	defer f.Close()

	if err != nil {
		return []FavoriteSheet{}
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		entry, err := toSheet(scanner.Text())
		if err != nil {
			panic(err)
		}

		list = append(list, entry)
	}

	return list
}

func toSheet(entry string) (FavoriteSheet, error) {
	separated := strings.Split(entry, "=")

	if len(separated) != 2 {
		return FavoriteSheet{}, errors.New("Config is corrupted")
	}

	return FavoriteSheet{
		Name: separated[0],
		ID:   separated[1],
	}, nil
}

func Get(name string) (string, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return "", err
	}

	buffer := bufio.NewScanner(f)

	for buffer.Scan() {
		entry, err := toSheet(buffer.Text())

		if err != nil {
			panic(err)
		}

		if entry.Name == name {
			return entry.ID, nil
		}
	}

	if err = buffer.Err(); err != nil {
		log.Fatal(err)
	}

	return "", errors.New("No entry found")
}
