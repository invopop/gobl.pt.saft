// Package test provides common functions for testing.
package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
)

// UpdateOut is a flag that can be set to update the example files in test/data and
// test/data/out.
var UpdateOut = flag.Bool("update", false, "Update the example files in test/data and test/data/out")

// LoadEnvelopes loads every test envelope in the test/data folder, sorted by name.
func LoadEnvelopes() ([]*gobl.Envelope, error) {
	files, err := filepath.Glob(Path("test", "data", "*.json"))
	if err != nil {
		return nil, err
	}

	sort.Strings(files) // Ensure consistent ordering

	envs := make([]*gobl.Envelope, len(files))
	for i, f := range files {
		envs[i], err = LoadEnvelope(filepath.Base(f))
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", filepath.Base(f), err)
		}
	}

	return envs, nil
}

// LoadEnvelope loads a test file from the test/data folder as a GOBL envelope and
// will rebuild it if necessary to ensure any changes are accounted for.
func LoadEnvelope(name string) (*gobl.Envelope, error) {
	path := Path("test", "data", name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err != nil {
		return nil, err
	}

	if err := env.Calculate(); err != nil {
		return nil, err
	}

	if err := env.Validate(); err != nil {
		return nil, err
	}

	if *UpdateOut {
		data, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
	}

	return env, nil
}

// SignEnvelope calculates and signs an envelope with a throwaway key.
func SignEnvelope(env *gobl.Envelope) error {
	if err := env.Calculate(); err != nil {
		return err
	}

	pkey := dsig.NewES256Key()
	return env.Sign(pkey)
}

// LoadFile reads a file relative to the project root.
func LoadFile(elements ...string) ([]byte, error) {
	return os.ReadFile(Path(elements...))
}

// Path joins the provided elements to the project root.
func Path(elements ...string) string {
	np := []string{RootPath()}
	np = append(np, elements...)
	return filepath.Join(np...)
}

// RootPath returns the project root, regardless of current working directory.
func RootPath() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}
