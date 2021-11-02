/* ASCII ART WEB YNOV INFORMATIQUE 2020 */
/* Copyright INGREMEAU, CLAMADIEU-THARAUD, MICHEL Romain 2020 */

package ascii

import (
	"bufio"
	"fmt"
	"os"
)

func Art(text, police, outputName string, output bool) (string, int) {
	if output {
		status := FileExist(outputName)
		if status != 0 {
			return "", status
		}
	}
	var wordList []string
	var status int
	etape := 0
	finalResult := ""

	for y, i := range text {
		if i == 10 {
			wordList = append(wordList, text[etape:y])
			etape = y + 1
		}
	}
	wordList = append(wordList, text[etape:])
	for i, ligne := range wordList {
		words := []rune(ligne)

		if len(words) > 0 && words[len(words)-1] == 13 {
			words = words[:len(words)-1]
		}

		for i := 0; i < len(words); i++ {
			if words[i] < 32 || words[i] > 126 {
				status = 400
				return "", status
			}
		}

		vals, status := openFiles(police)
		if status == 400 {
			return " ", status
		}
		result := fillArray(words, vals)

		if i == 0 {
			finalResult += printResult(result, words, output, outputName)
		} else {
			finalResult += "\n"
			finalResult += printResult(result, words, output, outputName)
		}
	}
	status = 200
	return finalResult, status
}

func FileExist(outputName string) int {
	// If file exist
	if len(outputName) > 4 && outputName[len(outputName)-4:] != ".txt" {
		outputName = outputName + ".txt"
	} else if outputName == ".txt" {
		return 400
	} else if len(outputName) <= 4 {
		outputName = outputName + ".txt"
	}
	_, err := os.Stat(outputName)
	if err == nil {
		if os.IsNotExist(err) {
			// Create file
			ff, err := os.Create(outputName)
			if err != nil && ff == nil {
				return 400
			}
		} else {
			//Delete file
			e := os.Remove(outputName)
			if e != nil {
				return 400
			}
			// Create file
			ff, err := os.Create(outputName)
			if err != nil && ff == nil {
				return 400
			}
		}
	}
	return 0
}

func openFiles(filename string) ([]string, int) {
	// Open file
	filename = filename + ".txt"
	file, err := os.Open(filename)
	if err != nil && file == nil {
		return nil, 400
	}
	var vals []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		vals = append(vals, scanner.Text())
	}
	return vals, 0
}

func fillArray(word []rune, vals []string) [][]string {
	//Fill array
	var lines = make([][]int, 8)
	for i := range lines {
		lines[i] = make([]int, len(word))
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < len(word); j++ {
			lines[i][j] = findLine(word, j, i)
		}
	}

	var result = make([][]string, 8)
	for i := range result {
		result[i] = make([]string, len(word))
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < len(word); j++ {
			result[i][j] = vals[lines[i][j]-1]
		}
	}
	return result
}

func findLine(word []rune, nb int, line int) int {
	return int((word[nb]-31)*8+(word[nb]-31)-6) + line - 1
}

func isInArray(array [][]int, value int, word []rune) (bool, int, int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < len(word); j++ {
			if int(array[i][j]) == value {
				return true, i, j
			}
		}
	}
	return false, 0, 0
}

func printResult(array [][]string, word []rune, mode bool, outputName string) string {
	// Mode == -1 -> Print in console
	// Else -> Print in file
	result := ""
	if !mode {
		for i := 0; i < 8; i++ {
			for j := 0; j < len(word); j++ {
				result += array[i][j]
				if j == len(word)-1 {
					result += "\n"
				}
			}
		}
	} else if mode {
		// Open file
		f, err := os.OpenFile(outputName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return result
		}

		//Write in file
		for i := 0; i < 8; i++ {
			line := ""
			for j := 0; j < len(word); j++ {
				line += array[i][j]
			}
			writer := bufio.NewWriter(f)
			fmt.Fprintln(writer, line)
			writer.Flush()
		}
	}
	return result
}

func wordLen(array [][]string) (int, int) {
	result := 0
	spaces := 0
	for i := 0; i < len(array[0]); i++ {
		result += len(array[0][i])
		if array[0][i] == "      " {
			spaces++
		}
	}
	return result, spaces
}
