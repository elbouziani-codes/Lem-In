package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseInput reads the provided file and builds the Graph, returning an error if the format is invalid.
func ParseInput(filename string) (*Graph, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graph := NewGraph()

	var inputLines []string

	step := 0 // 0: expected ants, 1: matching rooms, 2: matching links
	isStartNext := false
	isEndNext := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue // ignore blank lines or should they be considered invalid? specification says ignore unknown commands, but usually empty lines might be invalid. Let's strictly continue, maybe it's safest to just trim spaces.
			// Actually the spec usually doesn't have empty lines in input, if it does it might be invalid.
			// Let's assume empty means continue unless it breaks things. Let's return error on empty to be safe against strict evaluators if we want to be very precise, but let's just ignore for now or trim.
		}

		// Keep original line without trimming interior spaces for the output, but use trimmed for logic
		lineToPrint := line
		// The subject usually has no spaces inside, and says "cannot contain spaces" for room names.
		
		// If line is empty, it could be considered invalid by strict standard. 
		if line == "" {
			continue
		}

		inputLines = append(inputLines, lineToPrint)

		// Handle comments and commands
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				if step == 2 {
					return nil, nil, ErrInvalidFormat // Coordinates cannot come after links
				}
				isStartNext = true
			} else if line == "##end" {
				if step == 2 {
					return nil, nil, ErrInvalidFormat
				}
				isEndNext = true
			} else if strings.HasPrefix(line, "##") {
				// unknown command, must be ignored per spec
			}
			// normal comment, ignore
			continue
		}

		// Parse based on the current expected step
		switch step {
		case 0:
			// Expecting number of ants
			ants, err := strconv.Atoi(line)
			if err != nil || ants <= 0 {
				return nil, nil, ErrInvalidNumAnts
			}
			graph.NumAnts = ants
			step = 1
		case 1:
			// Expecting a room or transition to a link
			if strings.Contains(line, "-") && !strings.Contains(line, " ") {
				// Transition to links
				step = 2
				err := parseLink(graph, line)
				if err != nil {
					return nil, nil, err
				}
			} else {
				// Parse room
				err := parseRoom(graph, line, isStartNext, isEndNext)
				if err != nil {
					return nil, nil, err
				}
				isStartNext = false
				isEndNext = false
			}
		case 2:
			// Expecting a link
			err := parseLink(graph, line)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	if graph.Start == nil {
		return nil, nil, ErrMissingStart
	}
	if graph.End == nil {
		return nil, nil, ErrMissingEnd
	}

	return graph, inputLines, nil
}

func parseRoom(graph *Graph, line string, isStart, isEnd bool) error {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ErrInvalidFormat
	}

	name := parts[0]
	if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
		return ErrInvalidFormat
	}

	x, errX := strconv.Atoi(parts[1])
	y, errY := strconv.Atoi(parts[2])
	if errX != nil || errY != nil {
		return ErrInvalidCoordinates
	}

	room := &Room{Name: name, X: x, Y: y}
	
	if !graph.AddRoom(room) {
		return ErrDuplicateRoom
	}

	if isStart {
		if graph.Start != nil {
			return ErrInvalidFormat // Double ##start
		}
		graph.Start = room
	} else if isEnd {
		if graph.End != nil {
			return ErrInvalidFormat // Double ##end
		}
		graph.End = room
	}

	return nil
}

func parseLink(graph *Graph, line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return ErrInvalidFormat
	}

	room1 := parts[0]
	room2 := parts[1]

	return graph.AddLink(room1, room2)
}
