package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const nonExistentPath = "none/existent/path"
const dataDirectory = "data"
const tmpDataDirectory = "tmp_data"

var (
	pathWithExpectation = []struct {
		Path        string
		Expectation bool
	}{
		{
			nonExistentPath,
			false,
		},
		{
			filepath.Join(dataDirectory, "directory"),
			true,
		},
		{
			filepath.Join(dataDirectory, "file"),
			false,
		},
	}
	limitWithCount = []struct {
		Limit int
		Count int
	}{
		{2, 0},
		{1, 1},
		{0, 2},
	}
)

func removeTmpData() {
	err := os.RemoveAll(tmpDataDirectory)
	if err != nil {
		log.Fatalf("Could not remove %s directory", tmpDataDirectory)
	}
}

func recreateTmpData() {
	removeTmpData()
	cmd := exec.Command("cp", "-pR", dataDirectory, tmpDataDirectory)
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Could not copy %s directory", dataDirectory)
	}
}

func assertLimit(t *testing.T, expectedLimit int) {
	actualLimit := Limit()

	if expectedLimit != actualLimit {
		t.Errorf("Expected limit to be %v, got %v", expectedLimit, actualLimit)
	}
}

func assertPath(t *testing.T, expectedPath string) {
	actualPath := Path()

	if expectedPath != actualPath {
		t.Errorf("Expected path to be %s, got %s", expectedPath, actualPath)
	}
}

func TestLimit(t *testing.T) {
	assertLimit(t, 0)
}

func TestPath(t *testing.T) {
	assertPath(t, "")
}

func TestIsDirectory(t *testing.T) {
	for _, single := range pathWithExpectation {
		isDirectory := isDirectory(single.Path)

		if single.Expectation != isDirectory {
			t.Errorf("Expected %t for path %s, got %t", single.Expectation, single.Path, isDirectory)
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
	actualCount := len(directories)

	if expectedCount != actualCount {
		t.Errorf("Expected %v directories, got %v", expectedCount, actualCount)
	}
}

func TestHandleDirectory(t *testing.T) {
	defer removeTmpData()

	for _, single := range limitWithCount {
		recreateTmpData()
		count, _ := handleDirectory(tmpDataDirectory, single.Limit)

		if count != single.Count {
			t.Errorf("Expected %v directories to be removed, got %v", single.Count, count)
		}

		if single.Limit == 1 {
			files, err := ioutil.ReadDir(tmpDataDirectory)

			if err != nil {
				log.Fatal("Data directory should exist")
			}

			directories := filterDirectories(files)

			expectedName := "directory1"
			actualName := directories[0].Name()

			if actualName != expectedName {
				t.Errorf("Expected %s, got %s", expectedName, actualName)
			}
		}
	}
}

func TestRunCommand(t *testing.T) {
	defer removeTmpData()
	recreateTmpData()

	cases := []struct {
		Path          string
		Limit         int
		ExpectedText  string
		ExpectedError string
	}{
		{"none/existent/path", 2, "", "Not a directory or path \"none/existent/path\" does not exist!"},
		{tmpDataDirectory, 5, "", ""},
		{tmpDataDirectory, 1, "Removed a single directory.", ""},
		{tmpDataDirectory, 0, "Removed 2 directories.", ""},
	}

	for _, c := range cases {
		recreateTmpData()

		text, err := runCommand(c.Path, c.Limit)

		if err != nil {

		}

		if c.ExpectedText != text {
			t.Errorf("Wrong output for case: %v, expected %v, got: %v", c, c.ExpectedText, text)
		}

		if c.ExpectedError != "" && c.ExpectedError != err.Error() {
			t.Errorf("Wrong error for case: %v, expected %v, got: %v", c, c.ExpectedError, err)
		}
	}
}

func TestRealMain(t *testing.T) {
	defer removeTmpData()
	recreateTmpData()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		Name           string
		Args           []string
		ExpectedCode   int
		ExpectedOutput string
	}{
		{"limit set to 0 and tmp_data", []string{"-l", "0", tmpDataDirectory}, 0, "Removed 2 directories.\n"},
		{"limit set to default and tmp_data", []string{tmpDataDirectory}, 0, ""},
		{"limit set to default and wrong dir", []string{nonExistentPath}, 1, fmt.Sprintf("Not a directory or path \"%s\" does not exist!", nonExistentPath)},
		{"limit set to default and no dir", []string{}, 0, "Usage: cleanup path/to/dir -l 5\n  -l int\n    \tLimit of the latest directories to keep (default 5)\n"},
	}

	for _, c := range cases {
		flag.CommandLine = flag.NewFlagSet(c.Name, flag.ExitOnError)
		os.Args = os.Args[:1]
		os.Args = append(os.Args, c.Args...)

		var buf bytes.Buffer

		actualCode := realMain(&buf)
		if c.ExpectedCode != actualCode {
			t.Errorf("Wrong exit code for args: %v, expected: %v, got: %v", c.Args, c.ExpectedCode, actualCode)
		}

		actualOutput := buf.String()
		if c.ExpectedOutput != actualOutput {
			t.Errorf("Wrong output for args: %v, expected %v, got: %v", c.Args, c.ExpectedOutput, actualOutput)
		}
	}
}
