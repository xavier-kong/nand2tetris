package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readFile(filePath string) []string {
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

func handleLine(line string) string {
	if len(line) == 0 {
		return ""
	}

	return ""
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing args")
		return
	}

	filePath := os.Args[1]
	fileLines := readFile(filePath)

	res := []string

	for _, line := range fileLines {
		resLine := handleLine(line)

		if resLine != "" {
			append(res, resLine)
		}
	}

}
