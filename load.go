package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type destEnt struct {
	name1 string
	name2 string
	value int
}

func loadDestEntries(r io.Reader) (ents []destEnt, err error) {

	bufRdr := bufio.NewReader(r)
	var lineNo int

	for {
		var line string
		lineNo++
		if line, err = bufRdr.ReadString('\n'); len(line) == 0 && err == io.EOF {
			err = nil
			break
		} else if err == io.EOF {
			// ignore
		} else if err != nil {
			return nil, fmt.Errorf("input error at line %d: %s", lineNo, err)
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue // ignore empty lines
		}

		// first city name:
		index := strings.Index(line, "-")
		if index == -1 {
			return nil, fmt.Errorf("missing '-' at line %d", lineNo)
		}
		city1 := strings.TrimSpace(line[0:index])
		line = line[index+1:]

		// second city name:
		if index = strings.Index(line, ":"); index == -1 {
			return nil, fmt.Errorf("missing ':' at line %d", lineNo)
		}
		city2 := strings.TrimSpace(line[0:index])
		line = line[index+1:]

		// value:
		var value int64
		line = strings.TrimSpace(line)
		if value, err = strconv.ParseInt(line, 0, 0); err != nil {
			return nil, fmt.Errorf("invalid destination value %q at line %d", line, lineNo)
		}

		var ent destEnt
		ent.name1 = city1
		ent.name2 = city2
		ent.value = int(value)
		ents = append(ents, ent)
	}

	return
}

func mustLoadDestEntriesFromFile(filename string) (ents []destEnt) {
	if file, err := os.Open(filename); err != nil {
		panic(err)
	} else if ents, err = loadDestEntries(file); err != nil {
		panic(err)
	} else if err = file.Close(); err != nil {
		panic(err)
	}
	return
}

type routeEnt struct {
	name1 string
	name2 string
	dist  int
	color string
}

func loadRouteEntries(r io.Reader) (ents []routeEnt, err error) {

	bufRdr := bufio.NewReader(r)
	var lineNo int

	for {
		var line string
		lineNo++
		if line, err = bufRdr.ReadString('\n'); len(line) == 0 && err == io.EOF {
			err = nil
			break
		} else if err == io.EOF {
			// ignore
		} else if err != nil {
			return nil, fmt.Errorf("input error at line %d: %s", lineNo, err)
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue // ignore empty lines
		}

		// first city name:
		index := strings.Index(line, "-")
		if index == -1 {
			return nil, fmt.Errorf("missing '-' at line %d", lineNo)
		}
		city1 := strings.TrimSpace(line[0:index])
		line = line[index+1:]

		// second city name:
		if index = strings.Index(line, ":"); index == -1 {
			return nil, fmt.Errorf("missing ':' at line %d", lineNo)
		}
		city2 := strings.TrimSpace(line[0:index])
		line = line[index:]

		// route descriptions:
		line = strings.TrimSpace(line)
		for len(line) > 1 {
			var ent routeEnt
			line = strings.TrimSpace(line[1:]) // chop ':' or ',' and following whitespace
			if index = strings.IndexAny(line, " \t"); index == -1 {
				return nil, fmt.Errorf("missing route color at line %d", lineNo)
			}
			distText := strings.TrimSpace(line[:index])
			var dist int64
			if dist, err = strconv.ParseInt(distText, 0, 0); err != nil {
				return nil, fmt.Errorf("invalid route distance %q at line %d", distText, lineNo)
			}
			ent.dist = int(dist)
			line = strings.TrimSpace(line[index:])
			if index = strings.Index(line, ","); index == -1 {
				index = len(line)
			}
			ent.color = line[:index]
			ent.name1 = city1
			ent.name2 = city2
			ents = append(ents, ent)
			line = line[index:]
		}
	}

	return
}

func mustLoadRouteEntriesFromFile(filename string) (ents []routeEnt) {
	if file, err := os.Open(filename); err != nil {
		panic(err)
	} else if ents, err = loadRouteEntries(file); err != nil {
		panic(err)
	} else if err = file.Close(); err != nil {
		panic(err)
	}
	return
}

func mustLoadRouteEntriesFromString(s string) (ents []routeEnt) {
	var err error
	if ents, err = loadRouteEntries(strings.NewReader(s)); err != nil {
		panic(err)
	}
	return ents
}
