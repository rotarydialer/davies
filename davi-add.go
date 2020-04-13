package main

import (
	"os"
	"fmt"
	"io"
)

func main() {
	fmt.Println("Start.")
	files := os.Args[1:]

	fmt.Println(files)

	for _, element := range files {
		CopyToScripts(element)
	}

	fmt.Println("End.")
}

func CopyToScripts(filename string) {
	var err error
	destFilename := "scripts/" + filename // TODO: generate unique name

	if !fileExists(filename) {
		fmt.Println("ERROR: file not found", filename)
        os.Exit(1)
	}

	_, err = copyFile(filename, destFilename)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to copy '%s' to '%s'", filename, destFilename))
		fmt.Println(err)
	} else {
		fmt.Println(fmt.Sprintf("👍 '%s'", destFilename))
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
			return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
			return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
			return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
			return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}