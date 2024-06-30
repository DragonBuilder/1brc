package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type stat struct {
	Name  string
	Sum   float32
	Count int
	Min   float32
	Max   float32
}

func (s stat) Compare(o stat) int {
	return lexicalCompare(s.Name, o.Name)
	// if len(s.Name) == len(o.Name) {
	// 	if s.Name < o.Name {
	// 		return -1
	// 	}
	// 	if s.Name > o.Name {
	// 		return 1
	// 	}
	// 	return 0
	// }
	// defaultReturn := 0
	// if len(s.Name) < len(o.Name) {
	// 	o.Name = o.Name[:len(s.Name)]
	// 	defaultReturn = -1
	// } else {
	// 	s.Name = s.Name[:len(o.Name)]
	// 	defaultReturn = 1
	// }

	// if s.Name < o.Name {
	// 	return -1
	// }
	// if s.Name > o.Name {
	// 	return 1
	// }
	// return defaultReturn
}

func lexicalCompare(s1 string, s2 string) int {
	onEqualReturn := 0
	s1Len := len(s1)
	s2Len := len(s2)
	if s1Len != s2Len {
		if s1Len < s2Len {
			onEqualReturn = -1
			s2 = s2[:s1Len]
		} else {
			onEqualReturn = 1
			s1 = s1[:s2Len]
		}
	}
	if s1 < s2 {
		return -1
	}
	if s1 > s2 {
		return 1
	}
	return onEqualReturn
}

func (s *stat) Mean() float32 {
	return s.Sum / float32(s.Count)
}

func (s *stat) NewReading(temp float32) {
	s.Count += 1
	s.Sum += temp
	if temp < s.Min {
		s.Min = temp
	}

	if temp > s.Max {
		s.Max = temp
	}
}

func NewStat(name string, temp float32) *stat {
	return &stat{
		Name:  name,
		Sum:   temp,
		Count: 1,
		Min:   temp,
		Max:   temp,
	}
}

func (s *stat) OutputFormat() string {
	return fmt.Sprintf("%s=%.1f/%.1f/%.1f", s.Name, s.Min, s.Mean(), s.Max)
}

func (s *stat) Encoded() []byte {
	return []byte(fmt.Sprintf("%s;%f;%d;%f;%f", s.Name, s.Sum, s.Count, s.Min, s.Max))
}

func formatToOutputAlphabetically(stats map[string]*stat, order []string) []string {
	result := make([]string, 0)
	for _, name := range order {
		result = append(result, stats[name].OutputFormat())
	}
	return result
}

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

	fileReader, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Error while reading file %s : %v", inputFilePath, err)
		return
	}
	defer fileReader.Close()

	fileScanner := bufio.NewScanner(fileReader)

	if err := os.MkdirAll("out", os.ModePerm); err != nil {
		log.Fatalf("error creating out directory : %v", err)
	}

	cache := make(map[string]*stat)
	linesRead := 0
	var bst *Node
	for fileScanner.Scan() {
		linesRead += 1

		parsed := strings.Split(fileScanner.Text(), ";")
		name := parsed[0]
		temp, err := strconv.ParseFloat(parsed[1], 32)
		if err != nil {
			log.Fatalf("Error parsing float : %v", err)
		}

		if stationStat, ok := cache[name]; ok {
			stationStat.NewReading(float32(temp))
		} else {
			cache[name] = NewStat(name, float32(temp))
			if linesRead == 1 {
				bst = NewNode(cache[name])
			} else {
				BSTInsert(bst, cache[name])
			}

		}

	}

	stat_gen_time_taken := time.Since(start)

	// fmt.Printf("{%s}\n\n", strings.Join(formatToOutput(cache), ", "))
	fmt.Printf("{%s}\n\n", strings.Join(formatToOutputAlphabetically(cache, inorderTravese(bst)), ", "))

	total_time_taken := time.Since(start)
	fmt.Printf("Done.\n Took %s\n Calculations Took %s\n Total Stations: %d\n Total Lines: %d", total_time_taken.String(), stat_gen_time_taken.String(), len(cache), linesRead)
}

type Node struct {
	Value *stat
	Left  *Node
	Right *Node
}

func NewNode(val *stat) *Node {
	return &Node{Value: val}
}

func BSTInsert(root *Node, val *stat) {
	toInsert := NewNode(val)

	current := root
	for {

		comparison := toInsert.Value.Compare(*current.Value)
		if comparison < 0 {
			if current.Left != nil {
				current = current.Left
			} else {
				current.Left = toInsert
				break
			}
		} else if comparison > 0 {
			if current.Right != nil {
				current = current.Right
			} else {
				current.Right = toInsert
				break
			}
		} else {
			break
		}
	}
}

func inorderTravese(node *Node) []string {
	traversed := make([]string, 0)
	if node == nil {
		return traversed
	}
	traversed = append(traversed, inorderTravese(node.Left)...)
	traversed = append(traversed, node.Value.Name)
	traversed = append(traversed, inorderTravese(node.Right)...)

	return traversed
}
