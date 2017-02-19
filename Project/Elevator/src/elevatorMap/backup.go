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

	mapArray := NewCleanMap()

	for i := 0; i < def.Elevators; i++ {
		mapArray[def.IPs[i]].IP = stringMatrix[i*(3+def.Floors)][0]
		for j := 0; j < def.Floors; j++ {
			for k := 0; k < 3; k++ {
				mapArray[def.IPs[i]].Buttons[j][k], _ = strconv.Atoi(stringMatrix[i*(3+def.Floors)+1+j][k])
			}
		}
		mapArray[def.IPs[i]].Dir, _ = strconv.Atoi(stringMatrix[i*(3+def.Floors)+def.Floors+2][0])
		mapArray[def.IPs[i]].Pos, _ = strconv.Atoi(stringMatrix[i*(3+def.Floors)+def.Floors+1][0])
	}

	return mapArray

}

func WriteBackup(mapArray def.ElevMap) {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMatrix := [][]string{}

	for i := 0; i < def.Elevators; i++ {
		stringMatrix = append(stringMatrix, []string{mapArray[def.IPs[i]].IP})
		for j := 0; j < def.Floors; j++ {
			stringArray := []string{}
			for k := 0; k < 3; k++ {
				stringArray = append(stringArray, strconv.Itoa(mapArray[def.IPs[i]].Buttons[j][k]))
			}
			stringMatrix = append(stringMatrix, stringArray)
		}
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(mapArray[def.IPs[i]].Dir)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(mapArray[def.IPs[i]].Pos)})
	}
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMatrix)
	if err != nil {
		log.Fatal(err)
	}
}
