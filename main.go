package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ui-kreinhard/go-turing-i2c-cmdline/shell"
)

var nodeMapping map[int]string
var lockedNodes []int

func getOperation() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}
	return "", errors.New("no operation set")
}

func getNode() (int, error) {
	if len(os.Args) > 2 {
		nodeNumber, err := strconv.Atoi(os.Args[2])
		return nodeNumber, err
	}
	return -1, errors.New("no node given")
}

func logFatalIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func logCmdOutput(output string, err error) {
	if err != nil {
		log.Println(err, "cmd output", output)
	}
}

func main() {
	loadLockFile()
	initNodeMapping()
	operation, err := getOperation()
	logFatalIfErr(err)
	switch operation {
	case "turnOn":
		nodeNumber, err := getNode()
		logFatalIfErr(err)
		logCmdOutput(turnOnNode(nodeNumber))
	case "turnOff":
		nodeNumber, err := getNode()
		logFatalIfErr(err)
		logCmdOutput(turnOffNode(nodeNumber))
	case "turnOnAll":
		turnOnAllNodes()
	case "turnOffAll":
		turnOfAllNodes()
	default:
		log.Fatalln("No valid operation")
	}
}

func loadLockFile() error {
	lockedNodes = []int{}
	rawData, err := ioutil.ReadFile("/etc/lockedNodes")
	if err != nil {
		return err
	}
	data := strings.Split(string(rawData), ",")
	for _, nodeStr := range data {
		nodeStr = strings.TrimSpace(nodeStr)
		nodeInt, err := strconv.Atoi(nodeStr)
		if err == nil {
			lockedNodes = append(lockedNodes, int(nodeInt))
		} else {
			log.Println("cannot parse", nodeStr)
		}
	}
	log.Println("locked nodes are", lockedNodes, ". These nodes will be ignored for all operations")
	return nil
}

func isLocked(nodeNumber int) bool {
	for _, lockedNode := range lockedNodes {
		if lockedNode == nodeNumber {
			return true
		}
	}
	return false
}

func addNodeMapping(nodeNumber int, i2cAddress string) {
	if !isLocked(nodeNumber) {
		nodeMapping[nodeNumber] = i2cAddress
	}
}

func initNodeMapping() {
	nodeMapping = map[int]string{}
	addNodeMapping(1, "0x02")
	addNodeMapping(2, "0x04")
	addNodeMapping(3, "0x08")
	addNodeMapping(4, "0x10")
	addNodeMapping(5, "0x80")
	addNodeMapping(6, "0x40")
	addNodeMapping(7, "0x20")
}

func turnOffNode(nodeNumber int) (string, error) {
	nodeMask := nodeMapping[nodeNumber]
	return shell.Exec("i2cset", "-m", nodeMask, "-y", "1", "0x57", "0xf2", "0x00")
}

func turnOnNode(nodeNumber int) (string, error) {
	nodeMask := nodeMapping[nodeNumber]
	return shell.Exec("i2cset", "-m", nodeMask, "-y", "1", "0x57", "0xf2", "0xFF")
}

func turnOfAllNodes() {
	for nodeNumber, _ := range nodeMapping {
		logCmdOutput(turnOffNode(nodeNumber))
		time.Sleep(1 * time.Second)
	}
}

func turnOnAllNodes() {
	for nodeNumber, _ := range nodeMapping {
		logCmdOutput(turnOnNode(nodeNumber))
		time.Sleep(1 * time.Second)
	}
}
