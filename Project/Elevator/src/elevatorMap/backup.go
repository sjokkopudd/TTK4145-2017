package elevatorMap

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func readBackup() (mapArray [][]int) {
	backup, err := ioutil.ReadFile("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	tempReader := csv.NewReader(strings.NewReader(string(backup)))
	mapArray = [][]int{}

	for {
		stingLineFromCSV, err := tempReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		intLineFromCSV := []int{}

		for _, l := range stingLineFromCSV {
			tempInt, _ := strconv.Atoi(l)
			intLineFromCSV = append(intLineFromCSV, tempInt)
		}

		fmt.Println(intLineFromCSV)
		mapArray = append(mapArray, intLineFromCSV)

	}

	fmt.Println(mapArray)

	return mapArray
}

func writeBackup(mapArray [][]int) {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringArray := [][]string{}

	for i, _ := range mapArray {
		tempString := []string{}
		for _, l := range mapArray[i] {
			tempString = append(tempString, strconv.Itoa(l))
		}
		stringArray = append(stringArray, tempString)
	}

	fmt.Println(stringArray)

	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringArray)
	if err != nil {
		log.Fatal(err)
	}
}
