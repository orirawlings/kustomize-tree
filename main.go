package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func main() {
	for _, path := range os.Args[1:] {
		if err := tree(path, ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %q: %v\n", path, err)
		}
	}
}

func tree(path, indent string) error {
	var f io.ReadCloser
	var err error
	for _, name := range []string{
		"kustomization.yml",
		"kustomization.yaml",
		"Kustomization",
	} {
		f, err = os.Open(filepath.Join(path, name))
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("cannot open %q: %v", path, err)
	}
	defer f.Close()
	var data struct {
		Bases []string
	}
	if err := yaml.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("cannot parse %q: %v", path, err)
	}
	fmt.Println(path)
	for i, base := range data.Bases {
		add := "│  "
		if i == len(data.Bases)-1 {
			fmt.Printf(indent + "└──")
			add = "   "
		} else {
			fmt.Printf(indent + "├──")
		}

		if err := tree(filepath.Join(path, base), indent+add); err != nil {
			return err
		}
	}
	return nil
}
