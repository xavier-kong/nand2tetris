package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func handleComments(line string) string {
	if strings.HasPrefix(line, "//") { // full line comments
		return ""
	}

	parts := strings.Split(line, "//")
	return parts[0]
}

func handleAddressingInstruction(line string) string {
	addressString := strings.TrimPrefix(line, "@")
	addressInt, err := strconv.Atoi(addressString)

	if err != nil {
		panic(err)
	}

	return "0" + fmt.Sprintf("%015b", addressInt)
}

func handleLine(line string) string {
	line = handleComments(line)
	if len(line) == 0 {
		return ""
	}

	if strings.HasPrefix(line, "@") {
		line = handleAddressingInstruction(line)
	} else {
		line = handleComputeInstruction(line)
	}

	return line
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing args")
		return
	}

	filePath := os.Args[1]
	fileLines := readFile(filePath)

	var res []string

	for _, line := range fileLines {
		resLine := handleLine(line)

		if resLine != "" {
			res = append(res, resLine)
		}
	}

}
