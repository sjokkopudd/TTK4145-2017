package elevatorMap

import (
	"def"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func readBackup() def.ElevMap {
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
		newMap[e].Dir, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+1][0])
		newMap[e].Pos, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+2][0])
		newMap[e].Door, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+3][0])
		newMap[e].IsAlive, _ = strconv.Atoi(stringMatrix[e*(3+def.FLOORS)+def.FLOORS+4][0])
	}

	return newMap

}

func writeBackup(writeMap def.ElevMap) {
	backupFile, err := os.Create("src/elevatorMap/memory.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMatrix := [][]string{}

	for e := 0; e < def.ELEVATORS; e++ {
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[e].ID)})
		for f := 0; f < def.FLOORS; f++ {
			stringArray := []string{}
			for b := 0; b < def.BUTTONS; b++ {
				stringArray = append(stringArray, strconv.Itoa(writeMap[e].Buttons[f][b]))
			}
			stringMatrix = append(stringMatrix, stringArray)
		}
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[e].Dir)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[e].Pos)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[e].Door)})
		stringMatrix = append(stringMatrix, []string{strconv.Itoa(writeMap[e].IsAlive)})
	}
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMatrix)
	if err != nil {
		log.Fatal(err)
	}
}

func InitSoftwareBackup() {
	backupTicker := time.NewTicker(1 * time.Second)
	newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "make run")
	err := newBackup.Run()
	if err != nil {
		fmt.Println("Unable to spawn backup; you're on your own.")
		log.Fatal(err)
	}

	backupAdr, err := net.ResolveUDPAddr("udp", def.BACKUP_IP)
	if err != nil {
		log.Fatal(err)
	}

	backupConn, err := net.DialUDP("udp", nil, backupAdr)
	if err != nil {
		log.Fatal(err)
	}

	aliveMsg := true

	for {
		select {
		case <-backupTicker.C:
			jsonBuf, _ := json.Marshal(aliveMsg)
			backupConn.Write(jsonBuf)
		}

	}

}
