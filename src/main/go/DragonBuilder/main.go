package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

// const inputFilePath = "./measurements.txt"

func main() {
	slog.Info("Starting 1brc processing.")
	slog.Info(strings.Join(os.Args, ", "))
	if len(os.Args) < 2 {
		slog.Error("Provide a path to the file containing measurements.")
		return
	}
	inputFilePath := os.Args[1]
	if _, err := os.Stat(inputFilePath); errors.Is(err, os.ErrNotExist) {
		slog.Error("Input file not found.", "location", inputFilePath)
		return
	}

}
