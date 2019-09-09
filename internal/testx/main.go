package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	c := exec.Command("go", "env", "GOMOD")
	b, err := c.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	if len(b) == 0 {
		log.Fatal("we need modules")
	}

	root := strings.TrimSpace(string(b))
	root = filepath.Dir(root)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if err := os.Chdir(path); err != nil {
			return err
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_") {
			return filepath.SkipDir
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		infos, err := f.Readdir(-1)
		if err != nil {
			return err
		}

		var found bool
		for _, i := range infos {
			if strings.HasSuffix(i.Name(), "_test.go") {
				found = true
				break
			}
		}
		if !found {
			return nil
		}

		args := []string{"test"}
		args = append(args, os.Args[1:]...)
		c := exec.Command("go", args...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	})
	if err != nil {
		log.Fatal(err)
	}
}
