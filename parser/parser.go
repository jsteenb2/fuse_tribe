package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var (
	// Tracker is the collection of tracker data for each time frame
	Tracker []View
	// Heartbeat is the collection of heartbeat data and checks for connectivity
	Heartbeat []heartbeat
)

// ParseFile parses the eye-tribe output data, akin to JSON without the proper formatting
func ParseFile(filename string) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	parsedContent := strings.SplitAfter(string(f), "} ")
	parseTracker(parsedContent)
	parseHeartbeat(parsedContent)
}

func parseTracker(parsedContent []string) {
	for i := range parsedContent {
		var trackerData map[string]interface{}
		err := json.Unmarshal([]byte(parsedContent[i]), &trackerData)
		if err != nil {
			log.Fatal(err)
		}
		if trackerData["category"] == "tracker" {
			var result View
			mapstructure.Decode(trackerData, &result)
			Tracker = append(Tracker, result)
		}
	}
}

func parseHeartbeat(parsedContent []string) {
	for i := range parsedContent {
		var trackerData map[string]interface{}
		err := json.Unmarshal([]byte(parsedContent[i]), &trackerData)
		if err != nil {
			log.Fatal(err)
		}
		if trackerData["category"] == "heartbeat" {
			var result heartbeat
			mapstructure.Decode(trackerData, &result)
			Heartbeat = append(Heartbeat, result)
		}
	}
}

type heartbeat struct {
	Category   string
	Statuscode int
}

type View struct {
	Category   string
	Statuscode int
	Values     frame
}

type frame struct {
	Frame eyeData
}

type eyeData struct {
	Timestamp string
	Time      int
	fix       bool
	State     int
	Raw       xyCoordFloat
	Avg       xyCoordFloat
	Lefteye   eyeValue
	Righteye  eyeValue
}

type eyeValue struct {
	Raw     xyCoordFloat
	Avg     xyCoordFloat
	Psize   float64
	Pcenter xyCoordFloat
}

// type xyCoordInt struct {
// 	X int
// 	Y int
// }

type xyCoordFloat struct {
	X float64
	Y float64
}
