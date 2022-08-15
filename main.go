package main

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed template.md
var tepl string

var currentDate = "2006-01-02"

func ensureYaruDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	yaruHome := filepath.Join(home, ".yaru")
	err = os.MkdirAll(yaruHome, os.ModePerm)
	if os.IsExist(err) || err == nil {
		return yaruHome, nil
	}
	return "", err
}

func mainInternal() error {
	dir, err := ensureYaruDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "yaru-"+time.Now().Format(currentDate)+".md")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = file.WriteString(tepl)
		if err != nil {
			return err
		}
	}

	err = editFile(path)
	if err != nil {
		return err
	}

	return nil
}

func editFile(path string) error {
	command := getDefaultEditor() + " " + path
	cmd := exec.Command("sh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func getDefaultEditor() string {
	e := os.Getenv("EDITOR")
	if e != "" {
		return e
	}
	return "vim"
}

func main() {
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
}
