package main

import (
	"fmt"
	"os"
	messagepackage "proj3/messagePackage"
	"proj3/scheduler"
	"strconv"
	"sync"
	"time"
)

var (
	file *os.File
	mut  sync.Mutex
)

const usage = "Usage: runner mode threshold [number of threads]\n" +
	"mode     = (s) run sequentially, (p) process slices of each image in parallel \n" +
	"threshold = The threshold for intensity value of images for drawing contours \n " +
	"[number of threads] = Runs the parallel version of the program with the specified number of threads.\n"

func main() {
	// Parse the command-line arguments to get the number of consumers

	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}
	config := messagepackage.Config{Mode: "", ThreadCount: 0}
	config.Mode = os.Args[1]

	if len(os.Args) >= 4 {
		if config.Mode == "p" || config.Mode == "chunk" {
			threshold, _ := strconv.Atoi(os.Args[2])
			threads, _ := strconv.Atoi(os.Args[3])
			config.ThreadCount = threads
			config.Threshold = threshold
		} else {
			fmt.Println("Wrong Mode")
		}

	} else {
		threshold, _ := strconv.Atoi(os.Args[2])
		config.Threshold = threshold
	}

	start := time.Now()
	scheduler.Schedule(config)
	end := time.Since(start).Seconds()
	fmt.Printf("%.2f\n", end)
}
