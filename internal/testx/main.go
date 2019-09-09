package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	c := exec.Command("go", "env", "GOMOD")
	b, err := c.CombinedOutput()
	if err != nil {
		log.Fatal("GOMOD", err)
	}

	var root string

	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal("PWD", err)
		}
		root = pwd
	} else {
		root = filepath.Dir(string(b))
	}

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

		cfg := &packages.Config{
			Mode: packages.NeedFiles,
		}
		pkgs, err := packages.Load(cfg)
		if err != nil {
			log.Fatalf("error loading packages: %s, %s", path, err)
		}

		if len(pkgs) < 1 {
			return filepath.SkipDir
		}

		pkg := pkgs[0]
		if len(pkg.GoFiles) == 0 {
			return filepath.SkipDir
		}

		args := []string{"test"}
		args = append(args, os.Args[1:]...)
		c := exec.Command("go", args...)
		fmt.Printf(">  %s\n", path)
		fmt.Println("> ", strings.Join(c.Args, " "))
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	})
	if err != nil {
		log.Fatal(err)
	}
}
