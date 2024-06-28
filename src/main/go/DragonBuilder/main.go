package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"time"
)

// const inputFilePath = "./measurements.txt"

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		log.Fatalf("Input file path not given.\n\nExample: %s <path to file>", os.Args[0])
	}
	inputFilePath := os.Args[1]
	if _, err := os.Stat(inputFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("File not found: %s", inputFilePath)
		return
	}

	fileReader, err := os.OpenFile(inputFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while reading file %s : %v", inputFilePath, err)
		return
	}
	defer fileReader.Close()

	fileScanner := bufio.NewScanner(fileReader)
	for fileScanner.Scan() {
	}

	total_time_taken := time.Since(start)
	log.Printf("Done. Took %s", total_time_taken.String())

}
