package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var tempPath = filepath.Join(os.Getenv("TEMP"), "pipeline-files")

func main() {
	fmt.Println("start")
	start := time.Now()

	// generateFiles()
	processFiles()

	duration := time.Since(start)
	fmt.Println("done in", duration.Seconds(), "seconds")
}
