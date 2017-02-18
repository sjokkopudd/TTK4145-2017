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

func ReadBackup() [elevators]ElevatorInfo {
	backup, err := ioutil.ReadFile("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(strings.NewReader(string(backup)))
	csvReader.FieldsPerRecord = -1

	stringMatrix := [][]string{}

	for {
		csvLine, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("this error")
			log.Fatal(err)
		}
		stringMatrix = append(stringMatrix, csvLine)
	}

	mapArray := NewMap()

	for i := 0; i < elevators; i++ {
		mapArray[i].ID, _ = strconv.Atoi(stringMatrix[i*(3+floors)][0])
		for j := 0; j < floors; j++ {
			mapArray[i].Buttons[j].ButtonUp, _ = strconv.Atoi(stringMatrix[i*(3+floors)+1+j][0])
			mapArray[i].Buttons[j].ButtonUp, _ = strconv.Atoi(stringMatrix[i*(3+floors)+1+j][1])
			mapArray[i].Buttons[j].ButtonUp, _ = strconv.Atoi(stringMatrix[i*(3+floors)+1+j][2])
		}
		mapArray[i].Dir, _ = strconv.Atoi(stringMatrix[i*(3+floors)+floors+2][0])
		mapArray[i].Pos, _ = strconv.Atoi(stringMatrix[i*(3+floors)+floors+1][0])
	}

	return mapArray

}

func WriteBackup(mapArray [elevators]ElevatorInfo) {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMatrix := [][]string{}

	for i := 0; i < elevators; i++ {
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(mapArray[i].ID)})
		for j := 0; j < floors; j++ {
			stringArray := []string{}
			stringArray = append(stringArray, strconv.Itoa(mapArray[i].Buttons[j].ButtonUp))
			stringArray = append(stringArray, strconv.Itoa(mapArray[i].Buttons[j].ButtonDown))
			stringArray = append(stringArray, strconv.Itoa(mapArray[i].Buttons[j].ButtonPanel))
			stringMatrix = append(stringMatrix, stringArray)
		}
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(mapArray[i].Dir)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(mapArray[i].Pos)})
	}
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMatrix)
	if err != nil {
		log.Fatal(err)
	}
}
