package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var symbolTable = map[string]int{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"RO":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

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

func handleAddressingInstruction(line string) string {
	addressString := strings.TrimPrefix(line, "@")
	addressInt, err := strconv.Atoi(addressString)

	if err != nil {
		val, key := symbolTable[addressString]
		if key {
			addressInt = val
		}
	}

	return "0" + fmt.Sprintf("%015b", addressInt)
}

func findDest(destStr string) string {
	switch destStr {
	case "M":
		return "001"
	case "D":
		return "010"
	case "MD":
		return "011"
	case "A":
		return "100"
	case "AM":
		return "101"
	case "AD":
		return "110"
	case "AMD":
		return "111"
	default:
		return ""
	}
}

func findJump(jumpStr string) string {
	switch jumpStr {
	case "JGT":
		return "001"
	case "JEQ":
		return "010"
	case "JGE":
		return "011"
	case "JLT":
		return "100"
	case "JNE":
		return "101"
	case "JLE":
		return "110"
	case "JMP":
		return "111"
	default:
		return ""
	}
}

func findComp(compStr string) (comp string, a string) {
	switch compStr {
	case "0":
		return "101010", "0"
	case "1":
		return "111111", "0"
	case "-1":
		return "111010", "0"
	case "D":
		return "001100", "0"
	case "A":
		return "110000", "0"
	case "M":
		return "110000", "1"
	case "!D":
		return "001101", "0"
	case "!A":
		return "110001", "0"
	case "!M":
		return "110001", "1"
	case "-D":
		return "001111", "0"
	case "-A":
		return "110011", "0"
	case "-M":
		return "110011", "1"
	case "D+1":
		return "011111", "0"
	case "A+1":
		return "110111", "0"
	case "M+1":
		return "110111", "1"
	case "D-1":
		return "001110", "0"
	case "A-1":
		return "110010", "0"
	case "M-1":
		return "110010", "1"
	case "D+A":
		return "000010", "0"
	case "D+M":
		return "000010", "1"
	case "D-A":
		return "010011", "0"
	case "D-M":
		return "010011", "1"
	case "A-D":
		return "000111", "0"
	case "M-D":
		return "000111", "1"
	case "D&A":
		return "000000", "0"
	case "D&M":
		return "000000", "1"
	case "D|A":
		return "010101", "0"
	case "D|M":
		return "010101", "1"
	default:
		return "", ""
	}
}

func handleComputeInstruction(line string) string {
	hasEquals := strings.Contains(line, "=")
	hasCol := strings.Contains(line, ";")

	dest, comp, a, jump := "000", "", "", "000"

	if hasEquals {
		parts := strings.Split(line, "=")
		if hasCol { // dest=comp;jump
			compJump := strings.Split(parts[1], ";")
			parts = append(parts[:1], compJump...)
			dest, jump = findDest(parts[0]), findJump(parts[2])
			comp, a = findComp(parts[1])
		} else { // dest = comp
			dest = findDest(parts[0])
			comp, a = findComp(parts[1])
		}
	} else if hasCol { // comp; jump
		parts := strings.Split(line, ";")
		comp, a = findComp(parts[0])
		jump = findJump(parts[1])
	}

	return "111" + a + comp + dest + jump
}

func handleLine(line string) string {
	if len(line) == 0 {
		return ""
	}

	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "@") {
		line = handleAddressingInstruction(line)
	} else {
		line = handleComputeInstruction(line)
	}

	return line
}

func createOutput(res []string, path string) {
	filePath := strings.Replace(path, ".asm", ".hack", -1)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range res {
		fmt.Fprintln(w, line)
	}

	w.Flush()
}

func checkSymbolType(line string) string {
	var symbolType = ""
	matchLabel, _ := regexp.MatchString(`\(\S+\)`, line)
	matchVariable, _ := regexp.MatchString(`@[A-Za-z]+`, line)

	if matchLabel {
		symbolType = "label"
	} else if matchVariable {
		symbolType = "variable"
	}

	return symbolType
}

func handleCommentsAndSpace(lines []string) []string {
	var res []string

	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.Split(line, "//")
		if len(parts) == 2 {
			line = parts[0]
		}

		line = strings.TrimSpace(line)

		if len(line) != 0 {
			res = append(res, line)
		}

	}

	return res
}

func handleLabels(lines []string) []string {
	var res []string

	currLine := 0

	for _, line := range lines {
		symbolType := checkSymbolType(line)
		if symbolType == "label" {
			labelText := strings.ReplaceAll(line, "(", "")
			labelText = strings.ReplaceAll(labelText, ")", "")

			_, keyExists := symbolTable[labelText]

			if !keyExists {
				symbolTable[labelText] = currLine
			}

			continue
		}

		res = append(res, line)
		currLine += 1
	}

	return res
}

func handleVariables(lines []string) {
	currAdd := 16

	for _, line := range lines {
		symbolType := checkSymbolType(line)
		if symbolType == "variable" {
			varText := strings.ReplaceAll(line, "@", "")

			_, exists := symbolTable[varText]

			if !exists {
				symbolTable[varText] = currAdd
				currAdd += 1
			}
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing args")
		return
	}

	filePath := os.Args[1]
	fileLines := readFile(filePath)

	// handle comments and whitespace
	fileLines = handleCommentsAndSpace(fileLines)
	// handle labels
	fileLines = handleLabels(fileLines)
	// handle variables
	handleVariables(fileLines)

	var res []string

	for _, line := range fileLines {
		resLine := handleLine(line)

		if resLine != "" {
			res = append(res, resLine)
		}
	}

	createOutput(res, filePath)
}
