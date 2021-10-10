package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Equipment struct {
	TelescopeFocalLength int
	Mount                string
	Focuser              bool
	Guider               bool
	Rotator              bool
	Camera               string
}

type Target struct {
	Name            string
	SubExposureTime string
	SubExposureQty  int
	Filter          string
	Gain            string
	Offset          string
	Binning         string
	TimeStampStart  string
	TimeStampEnd    string
}

type Session struct {
	Targets   []Target
	Equipment Equipment
}

var focuserSequenceStart bool
var plateSolveStart bool
var meridianFlipStart bool

var equipment Equipment
var targets []Target
var tgt Target

func ProcessLogFile() Session {

	focuserSequenceStart = false
	equipment.TelescopeFocalLength = 0
	plateSolveStart = false
	meridianFlipStart = false

	file, err := os.Open(getLogFileName())
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		analyzeLogLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return Session{targets, equipment}
}

func analyzeLogLine(s string) {

	l := strings.Split(s, "|")
	if len(l) == 6 {
		d := strings.Split(l[5], " ")
		//Camera
		if strings.Contains(l[5], "Successfully connected Camera") {
			equipment.Camera = d[4] + d[5]
		}
		//has focuser
		if strings.Contains(l[5], "Successfully connected Focuser") {
			equipment.Focuser = true
		}
		//mount
		if strings.Contains(l[5], "Successfully connected Telescope") {
			equipment.Mount = d[4]
		}
		//guider
		if strings.Contains(l[5], "Guiding") {
			equipment.Guider = true
		}

		if strings.Contains(l[5], "Starting Category: Telescope, Item: Center") {
			plateSolveStart = true
		}
		if strings.Contains(l[5], "Finishing Category: Telescope, Item: Center") {
			plateSolveStart = false
		}

		//Telescope Focal Length
		if strings.Contains(l[5], "Platesolving with parameters") && equipment.TelescopeFocalLength == 0 {
			d = strings.SplitAfter(l[5], ":")
			fl := strings.SplitAfter(d[2], " ")
			var err error
			equipment.TelescopeFocalLength, err = strconv.Atoi(strings.TrimSpace(fl[1]))
			if err != nil {
				fmt.Printf(err.Error())
			}
		}

		//Target
		if strings.Contains(l[5], "Starting Category") && strings.Contains(l[5], "NINA.Sequencer.Container.DeepSkyObjectContainer") {
			d = strings.Split(l[5], ",")
			t := d[5]
			//targetSequenceStart = true
			tgt = Target{}
			tgt.Name = strings.TrimSuffix(strings.Split(t, ":")[1], "RA")
			tgt.TimeStampStart = l[0]
		}
		if strings.Contains(l[5], "Starting Category") && strings.Contains(l[5], "RunAutofocus") {
			focuserSequenceStart = true
		}
		if strings.Contains(l[5], "Finishing Category") && strings.Contains(l[5], "RunAutofocus") {
			focuserSequenceStart = false
		}
		if strings.Contains(l[5], "Meridian Flip - Recenter after meridian flip") {
			meridianFlipStart = true
		}
		if strings.Contains(l[5], "Meridian Flip - Settling scope") {
			meridianFlipStart = false
		}

		if strings.Contains(l[5], "Starting Exposure") && !focuserSequenceStart && !plateSolveStart && !meridianFlipStart {
			d = strings.Split(l[5], ";")
			fmt.Println("Starting Exposure")
			//Sub Time
			tgt.SubExposureTime = strings.Split(d[0], ":")[1]
			// Filter
			tgt.Filter = strings.Split(d[1], ":")[1]
			// Gain
			tgt.Gain = strings.Split(d[2], ":")[1]
			// Offset
			tgt.Offset = strings.Split(d[3], " ")[2]
			// Binning
			tgt.Binning = strings.Split(d[4], ":")[1]

			// Qty
			tgt.SubExposureQty++
		}
		if strings.Contains(l[5], "Finishing Category") && strings.Contains(l[5], "NINA.Sequencer.Container.DeepSkyObjectContainer") {
			//targetSequenceStart = false
			tgt.TimeStampEnd = l[0]
			targets = append(targets, tgt)
		}
	}
}

func getLogFileName() string {
	folder := "C:/Users/carlb/AppData/Local/NINA/Logs"
	file, err := os.Open(folder)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	fileList, _ := file.ReadDir(0)
	return folder + "/" + fileList[len(fileList)-1].Name()
}
