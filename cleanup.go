package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var (
	limit int
)

func Limit() int {
	return limit
}

func Path() string {
	return flag.Arg(0)
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func filterDirectories(files []fs.FileInfo) []fs.FileInfo {
	filtered := make([]fs.FileInfo, 0)

	for _, file := range files {
		if file.IsDir() {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

func handleDirectory(path string, limit int) (int, error) {
	count := 0
	removed := 0

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return 0, err
	}

	directories := filterDirectories(files)

	sort.Slice(directories, func(i, j int) bool {
		return directories[i].ModTime().Unix() > directories[j].ModTime().Unix()
	})

	for _, directory := range directories {
		count++

		if count > limit {
			err := os.RemoveAll(filepath.Join(path, directory.Name()))
			if err != nil {
				return removed, err
			}
			removed++
		}
	}

	return removed, nil
}

func runCommand(path string, limit int) (string, error) {
	if !isDirectory(path) {
		return "", errors.New(fmt.Sprintf("Not a directory or path \"%s\" does not exist!", path))
	}

	removed, err := handleDirectory(path, limit)

	if err != nil {
		return "", err
	}

	if removed > 0 {
		if removed == 1 {
			return fmt.Sprint("Removed a single directory."), nil
		} else {
			return fmt.Sprintf("Removed %d directories.", removed), nil
		}
	}

	return "", nil
}

func realMain(out io.Writer) int {
	flag.IntVar(&limit, "l", 5, "Limit of the latest directories to keep")

	flag.Usage = func() {
		_, err := fmt.Fprintln(out, "Usage: cleanup path/to/dir -l 5")
		if err != nil {
			return
		}
		flag.CommandLine.SetOutput(out)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return 0
	}

	text, err := runCommand(Path(), Limit())

	if err != nil {
		_, err := fmt.Fprint(out, err)
		if err != nil {
			return 1
		}
		return 1
	}

	if text != "" {
		_, err := fmt.Fprintln(out, text)
		if err != nil {
			return 1
		}
	}

	return 0
}

func main() {
	os.Exit(realMain(os.Stdout))
}
