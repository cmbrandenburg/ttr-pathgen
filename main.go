package main

import (
	"fmt"
	"os"
	"sort"
)

const (
	PROG_NAME = "ttr-pathgen"
)

func ePrintf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", PROG_NAME, msg)
}

func ePrintln(a ...interface{}) {
	new_args := []interface{}{fmt.Sprintf("%s:", PROG_NAME)}
	new_args = append(new_args, a...)
	fmt.Fprintln(os.Stderr, new_args...)
}

func makeDests() {

	// TODO: check for and remove duplicate destinations

	u := newUniv(mustLoadRouteEntriesFromFile("routes.dat"))
	dests := makeDestsEqualLikely(u)
	printDests(dests)
}

func printDests(dests []*dest) {
	for _, d := range dests {
		fmt.Printf("%q – %q : %d\n", d.city1.name, d.city2.name, d.value)
	}
}

func showDests() {
	u := newUniv(mustLoadRouteEntriesFromFile("routes.dat"))
	ents := mustLoadDestEntriesFromFile("destinations.dat")
	dests, err := newDestsFromDestEntries(u, ents)
	if err != nil {
		ePrintln(err)
		os.Exit(1)
	}
	printDests(dests)
}

func showRoutes() {
	u := newUniv(mustLoadRouteEntriesFromFile("routes.dat"))
	cities := u.allCitiesAlphabetical()
	for _, orig := range cities {
		numRoutes := 0
		numTgts := 0
		for _, routes := range orig.routes {
			numTgts++
			numRoutes += len(routes)
		}
		fmt.Printf("%q (%d, %d)\n", orig.name, numTgts, numRoutes)
		for tgt, routes := range orig.routes {
			for _, r := range routes {
				fmt.Printf("\t%q: %d %s\n", tgt.name, r.dist, r.color)
			}
		}
	}
}

func showShortestPaths() {
	u := newUniv(mustLoadRouteEntriesFromFile("routes.dat"))
	cities := u.allCitiesAlphabetical()
	for i, orig := range cities {
		for j, tgt := range cities {
			if i != j {
				pFewest := orig.fewestHops[tgt]
				pShortest := orig.shortestDist[tgt]
				var desc string
				if pFewest.equals(pShortest) {
					desc = fmt.Sprintf("%d hops, %d length", len(pFewest.routes), pFewest.dist)
				} else {
					desc = fmt.Sprintf("%d hops, %d length OR %d hops, %d length", len(pFewest.routes), pFewest.dist, len(pShortest.routes),
						pShortest.dist)
				}
				fmt.Printf("%q – %q: %s\n", orig.name, tgt.name, desc)
			}
		}
	}
}

func usage() {
	fmt.Println("usage:", PROG_NAME, "<command> <args>")
	fmt.Println()
	fmt.Println("Possible commands are:")
	fmt.Println("  make-dests")
	fmt.Println("  show-diff-len-paths")
	fmt.Println("  show-paths")
	fmt.Println("  show-routes")
	fmt.Println()
	os.Exit(1)
}

func main() {
	allCmds := map[string]func(){
		"make-dests":          makeDests,
		"show-dests":          showDests,
		"show-routes":         showRoutes,
		"show-shortest-paths": showShortestPaths,
	}
	if len(os.Args) <= 1 {
		var sortedCmds []string
		for cmd := range allCmds {
			sortedCmds = append(sortedCmds, cmd)
		}
		sort.Strings(sortedCmds)
		fmt.Println("usage:", PROG_NAME, "<command> <args>")
		fmt.Println()
		fmt.Println("Possible commands are:")
		for _, cmd := range sortedCmds {
			fmt.Printf("  %s\n", cmd)
		}
		fmt.Println()
		os.Exit(1)
	}
	cmd := os.Args[1]
	copy(os.Args[1:], os.Args[2:])
	os.Args = os.Args[:len(os.Args)-1]
	f := allCmds[cmd]
	if f != nil {
		f()
		return
	}
	ePrintf("invalid command %q", cmd)
}
