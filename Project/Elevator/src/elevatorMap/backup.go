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

func ReadBackup() [Elevators]ElevatorInfo {
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

	for i := 0; i < Elevators; i++ {
		mapArray[i].ID, _ = strconv.Atoi(stringMatrix[i*(3+Floors)][0])
		for j := 0; j < Floors; j++ {
			for k := 0; k < 3; k++ {
				mapArray[i].Buttons[j][k], _ = strconv.Atoi(stringMatrix[i*(3+Floors)+1+j][k])
			}
		}
		mapArray[i].Dir, _ = strconv.Atoi(stringMatrix[i*(3+Floors)+Floors+2][0])
		mapArray[i].Pos, _ = strconv.Atoi(stringMatrix[i*(3+Floors)+Floors+1][0])
	}

	return mapArray

}

func WriteBackup(mapArray [Elevators]ElevatorInfo) {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMatrix := [][]string{}

	for i := 0; i < Elevators; i++ {
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(mapArray[i].ID)})
		for j := 0; j < Floors; j++ {
			stringArray := []string{}
			for k := 0; k < 3; k++ {
				stringArray = append(stringArray, strconv.Itoa(mapArray[i].Buttons[j][k]))
			}
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
