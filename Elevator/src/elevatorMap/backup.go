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

func SoftwareBackup() {
	backupTicker := time.NewTicker(250 * time.Millisecond)
	newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "make run")
	err := newBackup.Run()
	if err != nil {
		fmt.Println("Unable to spawn backup; you're on your own.")
		return
	}

	backupAdr, err := net.ResolveUDPAddr("udp", def.BACKUP_IP)
	if err != nil {
		return
	}

	backupConn, err := net.DialUDP("udp", nil, backupAdr)
	if err != nil {
		return
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

func readBackup() ElevMap {
	backupFile, err := ioutil.ReadFile("src/elevatorMap/backup.txt")

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(strings.NewReader(string(backupFile)))
	csvReader.FieldsPerRecord = -1

	stringMap := [][]string{}

	for {
		csvLine, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		stringMap = append(stringMap, csvLine)
	}

	backupMap := *NewCleanElevMap()

	for e := 0; e < def.ELEVATORS; e++ {
		backupMap[e].ID, _ = strconv.Atoi(stringMap[e*(5+def.FLOORS)][0])
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				backupMap[e].Buttons[f][b], _ = strconv.Atoi(stringMap[e*(5+def.FLOORS)+1+f][b])
			}
		}
		backupMap[e].Direction, _ = strconv.Atoi(stringMap[e*(5+def.FLOORS)+def.FLOORS+1][0])
		backupMap[e].Position, _ = strconv.Atoi(stringMap[e*(5+def.FLOORS)+def.FLOORS+2][0])
		backupMap[e].Door, _ = strconv.Atoi(stringMap[e*(5+def.FLOORS)+def.FLOORS+3][0])
		backupMap[e].IsAlive, _ = strconv.Atoi(stringMap[e*(5+def.FLOORS)+def.FLOORS+4][0])
	}

	return backupMap

}

func writeBackup(backupMap ElevMap) {
	backupFile, err := os.Create("src/elevatorMap/backup.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer backupFile.Close()

	stringMap := [][]string{}

	for e := 0; e < def.ELEVATORS; e++ {
		stringMap = append(stringMap, []string{strconv.Itoa(backupMap[e].ID)})
		for f := 0; f < def.FLOORS; f++ {
			stringArray := []string{}
			for b := 0; b < def.BUTTONS; b++ {
				stringArray = append(stringArray, strconv.Itoa(backupMap[e].Buttons[f][b]))
			}
			stringMap = append(stringMap, stringArray)
		}
		stringMap = append(stringMap, []string{strconv.Itoa(backupMap[e].Direction)})
		stringMap = append(stringMap, []string{strconv.Itoa(backupMap[e].Position)})
		stringMap = append(stringMap, []string{strconv.Itoa(backupMap[e].Door)})
		stringMap = append(stringMap, []string{strconv.Itoa(backupMap[e].IsAlive)})
	}
	backupWriter := csv.NewWriter(backupFile)
	err = backupWriter.WriteAll(stringMap)
	if err != nil {
		log.Fatal(err)
	}
}
