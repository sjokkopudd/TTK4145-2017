package elevatorMap

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func readBackup() {
	backup, err := ioutil.ReadFile("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	tempReader := csv.NewReader(strings.NewReader(string(backup)))
	tempReader.Comment = '#'
	tempReader.TrimLeadingSpace = true
	lineNr := 0

	for {
		stingLineFromCSV, err := tempReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		mapArray[lineNr].firstFloorUp, _ = strconv.Atoi(stingLineFromCSV[0])
		mapArray[lineNr].secondFloorUp, _ = strconv.Atoi(stingLineFromCSV[1])
		mapArray[lineNr].secondFloorDown, _ = strconv.Atoi(stingLineFromCSV[2])
		mapArray[lineNr].thirdFloorUp, _ = strconv.Atoi(stingLineFromCSV[3])
		mapArray[lineNr].thirdFloorDown, _ = strconv.Atoi(stingLineFromCSV[4])
		mapArray[lineNr].fourthFloorDown, _ = strconv.Atoi(stingLineFromCSV[5])
		for i := 0; i < 3; i++ {
			mapArray[lineNr].elevatorPos[i], _ = strconv.Atoi(stingLineFromCSV[6+i])
			mapArray[lineNr].elevatorDir[i], _ = strconv.Atoi(stingLineFromCSV[9+i])
		}

		lineNr++
	}
	fmt.Println(mapArray)
}

func writeBackup() {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMatrix := [][]string{}

	for i := 0; i < 3; i++ {
		stringArray := []string{}

		stringArray = append(stringArray, strconv.Itoa(mapArray[i].firstFloorUp))
		stringArray = append(stringArray, strconv.Itoa(mapArray[i].secondFloorUp))
		stringArray = append(stringArray, strconv.Itoa(mapArray[i].secondFloorDown))
		stringArray = append(stringArray, strconv.Itoa(mapArray[i].thirdFloorUp))
		stringArray = append(stringArray, strconv.Itoa(mapArray[i].thirdFloorDown))
		stringArray = append(stringArray, strconv.Itoa(mapArray[i].fourthFloorDown))

		for j := 0; j < 3; j++ {
			stringArray = append(stringArray, strconv.Itoa(mapArray[i].elevatorPos[j]))

		}
		for j := 0; j < 3; j++ {
			stringArray = append(stringArray, strconv.Itoa(mapArray[i].elevatorDir[j]))
		}

		stringMatrix = append(stringMatrix, stringArray)

	}

	fmt.Println(stringMatrix)
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMatrix)
	if err != nil {
		log.Fatal(err)
	}

	commentWriter := bufio.NewWriter(backupFile)

	fmt.Fprintln(commentWriter, "#1u   2u   2d   3u   3d   4d | p1   p2   p3 | d1   d2   d3 ")
	commentWriter.Flush()

}
