package main

import (
	"encoding/json"
	"fmt"
	"go-event-processor/cmd"
	"go-event-processor/processor"
	"go-event-processor/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/stoicperlman/fls"
)


func main() {

	var ch = make(chan []byte)
	var wg sync.WaitGroup

	durations := map[string]types.DeliveryTimes{}

	var file_name string
	var window_size int

	cmd := cmd.NewProcessCmd()
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	file_name, err = cmd.Flags().GetString("file_name")
	if err != nil {
		log.Fatal(err)
	}
	window_size, err = cmd.Flags().GetInt("window_size")
	if err != nil {
		log.Fatal(err)
	}

	go processor.Process(&wg, ch, durations)

	osFile, err := os.Open(file_name)
	if err != nil {
		fmt.Printf("ERROR: failed to open input file %v\n", err)
		return
	}
	defer osFile.Close()

	now := time.Now().In(time.UTC)
	minTimestamp := now.Add(-time.Duration(window_size) * time.Minute)

	// Search lines starting from the bottom
	flsFile := fls.LineFile(osFile)

	previous, _ := flsFile.SeekLine(0, io.SeekEnd)
	byteSlice := make([]byte, previous)
	numBytes, _ := osFile.Read(byteSlice)
	line := string(byteSlice[:numBytes])
	if !processor.TimestampIsValid(minTimestamp, line) {
		return
	}
	// Initialize workers to get aggregations
	wg.Add(1)
	ch <- []byte(line)
	i := 0
	for {
		current, _ := flsFile.SeekLine(int64(-i), io.SeekEnd)
		byteSlice = make([]byte, previous-current)
		numBytes, _ = osFile.Read(byteSlice)
		line = string(byteSlice[:numBytes])
		if current != previous && processor.TimestampIsValid(minTimestamp, line) {
			wg.Add(1)
			ch <- []byte(line)
		}

		if current == 0 {
			break // reached end of file
		}

		i += 1
		previous = current
	}

	// Waits for all workers to finish to complete aggregations
	close(ch)
	wg.Wait()

	// Calculates averages
	averages := processor.Av(durations)

	// Outputs pretty JSON
	file, err := json.MarshalIndent(averages, "", " ")
	if err != nil {
		fmt.Printf("ERROR: failed to marshal indent output %v\n", err)
		return
	}
	err = ioutil.WriteFile("average_delivery_times.json", file, 0644)
	if err != nil {
		fmt.Printf("ERROR: failed write output to file %v\n", err)
		return
	}
}