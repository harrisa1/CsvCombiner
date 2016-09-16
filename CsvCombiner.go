package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const DIR = "."
const CSV = ".csv"
const OUTPUTFILE = "results.csv"
const HEADER = 0
const NOTFOUND = -1

var HEADERSTOIGNORE = []string{"State", "Total Interactions to Place", "Outstanding", "Active"}

var files []os.FileInfo
var currentFileName string
var currentFile os.File

var masterCsv [][]string

var currentCsv [][]string
var currentRow int
var columnsMissingFromCurrent []int
var extraColumnsInCurrent []int

func main() {
	getFiles()
	readFiles()
	writeOutputFile()

	fmt.Println("Press enter to exit")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getFiles() {
	var err error
	files, err = ioutil.ReadDir(DIR)
	if err != nil {
		fmt.Println(err)
	}
}

func readFiles() {
	for _, file := range files {
		currentFileName = file.Name()
		readFile()
	}
}

func readFile() {
	if currentFileName[len(currentFileName)-len(CSV):] == CSV {
		readCsvFile()
	}
}

func readCsvFile() {
	var err error
	fmt.Println("Reading file: " + currentFileName)
	currentFile, err := os.Open(currentFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer currentFile.Close()
	resetCurrent()
	reader := csv.NewReader(currentFile)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		currentCsv = append(currentCsv, record)
	}
	appendToMaster()
}

func resetCurrent() {
	currentCsv = [][]string{}
	currentRow = 0
	columnsMissingFromCurrent = []int{}
	extraColumnsInCurrent = []int{}
}

func appendToMaster() {
	diffColumns()
	mergeHeaders()
	appendRows()
	fixExistingRows()
}

func diffColumns() {
	diffMissingColumns()
	diffExtraColumns()
}

func diffMissingColumns() {
	if len(masterCsv) > 0 {
		for index, item := range masterCsv[HEADER] {
			if (len(currentCsv) > 0 && contains(currentCsv[HEADER], item)) == false {
				columnsMissingFromCurrent = append(columnsMissingFromCurrent, index)
			}
		}
	}
}

func diffExtraColumns() {
	if len(currentCsv) > 0 {
		for index, item := range currentCsv[HEADER] {
			if (len(masterCsv) > 0 && contains(masterCsv[HEADER], item)) == false {
				extraColumnsInCurrent = append(extraColumnsInCurrent, index)
			}
		}
	}
}

func mergeHeaders() {
	if len(masterCsv) > 0 {
		for _, item := range extraColumnsInCurrent {
			masterCsv[HEADER] = append(masterCsv[HEADER], currentCsv[HEADER][item])
		}
	} else {
		masterCsv = append(masterCsv, currentCsv[HEADER])
	}
}

func appendRows() {
	for currentRow = range currentCsv {
		if currentRow > 0 && contains(HEADERSTOIGNORE, currentCsv[currentRow][HEADER]) == false {
			if rowExists() {
				updateRow()
			} else {
				insertRow()
			}
		}
	}
}

func rowExists() bool {
	for _, row := range masterCsv {
		if row[HEADER] == currentCsv[currentRow][HEADER] {
			return true
		}
	}
	return false
}

func updateRow() {
	masterRow := getMasterRow()
	currentIndex := 1
	for i := 1; i < len(masterCsv[masterRow]); i++ {
		if containsInt(columnsMissingFromCurrent, i) == false {
			previousVal, _ := strconv.Atoi(masterCsv[masterRow][i])
			newVal, _ := strconv.Atoi(currentCsv[currentRow][currentIndex])
			masterCsv[masterRow][i] = strconv.Itoa(previousVal + newVal)
			currentIndex++
		}
	}

	for _, item := range extraColumnsInCurrent {
		masterCsv[masterRow] = append(masterCsv[masterRow], currentCsv[currentRow][item])
	}
}

func getMasterRow() int {
	for index, row := range masterCsv {
		if row[HEADER] == currentCsv[currentRow][HEADER] {
			return index
		}
	}
	fmt.Println("!!!!!!!!!!!!!!! Row not found to update !!!!!!!!!!!!!!!!")
	return NOTFOUND
}

func insertRow() {
	newRow := make([]string, len(masterCsv[HEADER]))
	currentIndex := 0
	for i := 0; i < len(masterCsv[HEADER]); i++ {
		if containsInt(columnsMissingFromCurrent, i) {
			newRow[i] = "0"
		} else {
			newRow[i] = currentCsv[currentRow][currentIndex]
			currentIndex++
		}
	}
	masterCsv = append(masterCsv, newRow)
}

func fixExistingRows(){
	cols := len(masterCsv[HEADER])
	for index, row := range masterCsv {
		tempRow := make([]string, cols)
		for i := 0; i < cols; i++ {
			if (i < len(row)){
				tempRow[i] = row[i]
			} else {
				tempRow[i] = "0"
			}
		}
		masterCsv[index] = tempRow
	}
}

func writeOutputFile() {
	fmt.Println("Writing output file: " + OUTPUTFILE)
	outputFile, err := os.Create(OUTPUTFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)
	writer.WriteAll(masterCsv)
}

func contains(array []string, s string) bool {
	return getIndex(array, s) > NOTFOUND
}

func getIndex(array []string, s string) int {
	for index, item := range array {
		if item == s {
			return index
		}
	}
	return NOTFOUND
}

func containsInt(array []int, i int) bool {
	for _, item := range array {
		if item == i {
			return true
		}
	}
	return false
}
