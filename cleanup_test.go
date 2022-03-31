package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type PathWithExpectation struct {
	path        string
	expectation bool
	error       string
}

type LimitWithCount struct {
	limit int
	count int
}

const dataDirectory = "data"
const tmpDataDirectory = "tmp_data"

var (
	data = []PathWithExpectation{
		{
			path:        "none/existent/path",
			expectation: false,
			error:       "",
		},
		{
			path:        filepath.Join(dataDirectory, "directory"),
			expectation: true,
			error:       "",
		},
		{
			path:        filepath.Join(dataDirectory, "file"),
			expectation: false,
			error:       "not a valid directory",
		},
	}
)

func removeTmpData() {
	err := os.RemoveAll(tmpDataDirectory)
	if err != nil {
		log.Fatal("Could not remove tmp_data directory")
	}
}

func recreateTmpData() {
	removeTmpData()
	cmd := exec.Command("cp", "-pR", dataDirectory, tmpDataDirectory)
	err := cmd.Run()

	if err != nil {
		log.Fatal("Could not copy data directory")
	}
}

func TestIsDirectory(t *testing.T) {
	for _, pwe := range data {
		isDirectory, _ := isDirectory(pwe.path)

		if isDirectory != pwe.expectation {
			t.Errorf("Failed! Expected %t for path %s, but got %t", pwe.expectation, pwe.path, isDirectory)
		}
	}
}

func TestEnsureDirectory(t *testing.T) {
	for _, pwe := range data {
		err := ensureDirectory(pwe.path)

		if pwe.expectation {
			if err != nil {
				t.Errorf("Failed! Got error %v", err)
			}
		} else {
			if err == nil {
				t.Errorf("Failed! Did not get an expected error")
			} else if len(pwe.error) != 0 {
				if err.Error() != pwe.error {
					t.Errorf("Failed! Expected error %s but got %v", pwe.error, err)
				}
			}
		}
	}
}

func TestFilterDirectories(t *testing.T) {
	files, err := ioutil.ReadDir(dataDirectory)

	if err != nil {
		log.Fatal("Data directory should exist")
	}

	directories := filterDirectories(files)

	expectedCount := 2
	gotCount := len(directories)

	if gotCount != expectedCount {
		t.Errorf("Failed! got %v directories and expected %v", gotCount, expectedCount)
	}
}

func TestHandleDirectory(t *testing.T) {
	defer removeTmpData()

	data := []LimitWithCount{
		{2, 0},
		{1, 1},
		{0, 2},
	}
	for _, single := range data {
		recreateTmpData()
		count, _ := handleDirectory(tmpDataDirectory, single.limit)

		if count != single.count {
			t.Errorf("Failed! Got %d directories removed and expected %d", count, single.count)
		}

		if single.limit == 1 {
			files, err := ioutil.ReadDir(tmpDataDirectory)

			if err != nil {
				log.Fatal("Data directory should exist")
			}

			directories := filterDirectories(files)

			if directories[0].Name() != "directory1" {
				t.Errorf("Failed! Got %s and expected %s", directories[0].Name(), "directory1")
			}
		}
	}
}
