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

	newMap := NewCleanMap()

	for i := 0; i < def.ELEVATORS; i++ {
		newMap[def.IPs[i]].IP = stringMatrix[i*(3+def.FLOORS)][0]
		for j := 0; j < def.FLOORS; j++ {
			for k := 0; k < 3; k++ {
				newMap[def.IPs[i]].Buttons[j][k], _ = strconv.Atoi(stringMatrix[i*(3+def.FLOORS)+1+j][k])
			}
		}
		newMap[def.IPs[i]].Dir, _ = strconv.Atoi(stringMatrix[i*(3+def.FLOORS)+def.FLOORS+1][0])
		newMap[def.IPs[i]].Pos, _ = strconv.Atoi(stringMatrix[i*(3+def.FLOORS)+def.FLOORS+2][0])
		newMap[def.IPs[i]].Door, _ = strconv.Atoi(stringMatrix[i*(3+def.FLOORS)+def.FLOORS+3][0])
	}

	return newMap

}

func WriteBackup(writeMap def.ElevMap) {
	
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	stringMatrix := [][]string{}

	for i := 0; i < def.ELEVATORS; i++ {
		stringMatrix = append(stringMatrix, []string{writeMap[def.IPs[i]].IP})
		for j := 0; j < def.FLOORS; j++ {
			stringArray := []string{}
			for k := 0; k < 3; k++ {
				stringArray = append(stringArray, strconv.Itoa(writeMap[def.IPs[i]].Buttons[j][k]))
			}
			stringMatrix = append(stringMatrix, stringArray)
		}
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[def.IPs[i]].Dir)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[def.IPs[i]].Pos)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[def.IPs[i]].Door)})
	}
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMatrix)
	if err != nil {
		fmt.Println("Stuck")
		log.Fatal(err)
	}

	err = backupFile.Close()

	if err != nil {
		fmt.Println("Could not close file")
		log.Fatal(err)
	}

}
