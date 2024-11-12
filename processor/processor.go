package processor

import (
	"encoding/json"
	"fmt"
	"go-event-processor/types"
	"sort"
	"strings"
	"sync"
	"time"
)

// aggregate groups delivery durations by timestamps's minute
func aggregate(event *types.Event, durations map[string]types.DeliveryTimes) {

	if event.EventName == "translation_delivered" {
		timestampSlice := strings.Split(event.Timestamp, ":")[:2]
		timestamp := timestampSlice[0] + ":" + timestampSlice[1] + ":00"
		if (durations[timestamp] == types.DeliveryTimes{}) {
			durations[timestamp] = types.DeliveryTimes{
				DurationSum: 0,
				Count:       0,
			}
		}
		currentDuration := durations[timestamp]
		currentDuration.DurationSum += event.Duration
		currentDuration.Count += 1
		var mu sync.Mutex
		mu.Lock() // Locks timestamp duration update until complete
		durations[timestamp] = currentDuration
		mu.Unlock()
	}
}

func Process(wg *sync.WaitGroup, ch chan []byte, durations map[string]types.DeliveryTimes) {

	for buf := range ch {
		var mu sync.Mutex
		var event types.Event
		mu.Lock() // Locks event struct unmarshal until complete
		err := json.Unmarshal(buf, &event)
		mu.Unlock()
		if err != nil {
			fmt.Println("json:", err)
			wg.Done()
			continue
		}
		aggregate(&event, durations)
		wg.Done()
	}
}

func TimestampIsValid(inputTimestamp time.Time, line string) bool {

	lineTimestampFragment := strings.Split(line, ",")[0]
	lineTimestampStr := strings.Split(strings.Split(strings.Split(lineTimestampFragment, "\"timestamp\": \"")[1], "\"")[0], "\"")[0]
	lineTimestamp, err := time.Parse(types.Timelayout, lineTimestampStr)

	if err != nil {
		fmt.Printf("ERROR: failed to parse line timestamp %v", err)
		return false
	}
	if lineTimestamp.Unix() >= inputTimestamp.Unix() {
		return true
	}
	return false
}

func Av(durations map[string]types.DeliveryTimes) []map[string]interface{} {

	averages := []map[string]interface{}{}

	// sort durations by timestamp
	var timestamps []string
	for timestamp := range durations {
		timestamps = append(timestamps, timestamp)
	}
	sort.Strings(timestamps)

	// calculate average
	for _, timestamp := range timestamps {
		duration := durations[timestamp]
		durations := float64(duration.DurationSum)
		count := float64(duration.Count)
		averageValue := durations / count
		average := map[string]interface{}{
			"date":                  timestamp,
			"average_delivery_time": averageValue,
		}
		averages = append(averages, average)
	}
	return averages
}