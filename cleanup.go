package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func ensureDirectory(path string) error {
	isDirectory, err := isDirectory(path)

	if err != nil {
		return err
	}

	if !isDirectory {
		return errors.New("not a valid directory")
	}

	return nil
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

func main() {
	var (
		limit int
		path  string
	)

	flag.IntVar(&limit, "l", 5, "Limit of latest directories to keep. Default value is 5.")

	flag.Usage = func() {
		fmt.Println("Usage: cleanup /path/to/dir -l 5")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	path = os.Args[1]

	if strings.TrimSpace(path) == "" {
		log.Fatal("Directory path is required!")
	}

	err := ensureDirectory(path)

	if err != nil {
		log.Fatal(err)
	}

	removed, err := handleDirectory(path, limit)

	if err != nil {
		log.Fatal(err)
	}

	if removed > 0 {
		if removed == 0 {
			fmt.Println("Removed a single directory.")
		} else {
			fmt.Println("Removed", removed, "directories.")
		}
	}
}
