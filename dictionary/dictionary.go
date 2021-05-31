package dictionary

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Tag struct {
	Name      string
	Tags      string
	UpdatedAt string
}

type Dictionary interface {
	Init() error
	SetIndex(string) error
	Create()
	Get()
	Synoyms()
	GetSynoym()
	Update() error
	Delete(string) error
}

func appendToDictionary(file string, keyword string, text string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	val, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	ss := strings.Split(string(val), "\n")
	content := ""
	isExist := false
	for _, line := range ss {
		if strings.HasPrefix(line, keyword) {
			line = line + "," + text
			isExist = true
		}
		content += line
	}

	if isExist == false {
		content += keyword + "," + text
	}

	err = writeToFile(file, content)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func writeToFile(file, content string) error {
	err := ioutil.WriteFile(file, []byte(content), 0600)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
