package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed template.md
var tepl string

func NewRootCommnad() *cobra.Command {
	cmd := cobra.Command{
		Use:   "yaru",
		Short: "yaru is daily note app",
		Long:  "yaru is daily note app",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root()
		},
	}
	cmd.AddCommand(NewListCommand())
	return &cmd
}

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

func root() error {
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

	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return clipboard.WriteAll(string(bs))
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

func NewListCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "list notes",
		Long:  "list notest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list()
		},
	}
	return &cmd
}

func list() error {
	path, err := ensureYaruDir()
	if err != nil {
		return err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) > 0
	})
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, e := range files {
		fmt.Fprintln(w, e.Name())
	}
	return nil
}
