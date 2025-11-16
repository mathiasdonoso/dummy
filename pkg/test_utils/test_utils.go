package testutils

import (
	"os"
	"testing"
)

func MustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading file %s: %v", path, err)
	}
	return data
}
