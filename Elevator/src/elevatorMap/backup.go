package elevatorMap

import (
	"def"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadBackup() def.ElevMap {
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

	newMap := *def.NewCleanElevMap()

	for e := 0; e < def.ELEVATORS; e++ {
		newMap[e].ID, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)][0])
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				newMap[e].Buttons[f][b], _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+1+f][b])
			}
		}
		newMap[e].State, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+1][0])
		newMap[e].Dir, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+2][0])
		newMap[e].Pos, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+3][0])
		newMap[e].Door, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+4][0])
		newMap[e].IsAlive, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+5][0])
	}

	return newMap

}

func WriteBackup(localMap def.ElevMap) {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMatrix := [][]string{}

	for e := 0; e < def.ELEVATORS; e++ {
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(localMap[e].ID)})
		for f := 0; f < def.FLOORS; f++ {
			stringArray := []string{}
			for b := 0; b < def.BUTTONS; b++ {
				stringArray = append(stringArray, strconv.Itoa(localMap[e].Buttons[f][b]))
			}
			stringMatrix = append(stringMatrix, stringArray)
		}
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(localMap[e].State)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(localMap[e].Dir)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(localMap[e].Pos)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(localMap[e].Door)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(localMap[e].IsAlive)})
	}
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMatrix)
	if err != nil {
		log.Fatal(err)
	}
}
