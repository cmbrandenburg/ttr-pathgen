package main

import (
	"fmt"
	"sort"
	"sync"
)

// TERMINOLOGY
//
// Route: A single connecting path between two adjacent cities. A route
// comprises a distance and color.
//
// Hop: A pair comprising a city and a route. The typical meaning of a hop is a
// route that leads to the given city.
//
// Path: A path is a origin city plus a continuous sequence of hops connecting
// two (possibly non-adjacent) cities. A path is (1) a sequence of N cities,
// where N >= 1; (2) a sequence of N-1 routes, with the routes being the
// connectors between the cities; and (3) the total distance of all routes.
//
// Destination: A pair of two cities connected by a path. A destination
// comprises the two cities and a value measured in points.
//

type pathComparer func(*path, *path) int

type city struct {
	name         string
	routes       map[*city][]*route
	fewestHops   map[*city]*path
	shortestDist map[*city]*path
}

func newCity(name string) *city {
	return &city{
		name:         name,
		routes:       make(map[*city][]*route),
		fewestHops:   make(map[*city]*path),
		shortestDist: make(map[*city]*path),
	}
}

func newCityMapFromRouteEntries(ents []routeEnt) (m map[string]*city) {
	m = make(map[string]*city)
	for _, ent := range ents {
		// Step 1: Populate the map with both cities, creating an empty city for any
		// that isn't already created.
		c1 := m[ent.name1]
		if c1 == nil {
			c1 = newCity(ent.name1)
			m[ent.name1] = c1
		}
		c2 := m[ent.name2]
		if c2 == nil {
			c2 = newCity(ent.name2)
			m[ent.name2] = c2
		}
		// Step 2: Populate both cities with the route.
		r := newRoute(ent.dist, ent.color)
		c1.routes[c2] = append(c1.routes[c2], r)
		c2.routes[c1] = append(c2.routes[c1], r)
	}
	return
}

func (c *city) populatePaths() {

	// TODO: test

	compFewestHops := func(p1, p2 *path) int {
		if len(p1.routes) < len(p2.routes) || (len(p1.routes) == len(p2.routes) && p1.dist < p2.dist) {
			return 1
		} else if len(p1.routes) == len(p2.routes) && p1.dist == p2.dist {
			return 0
		}
		return -1
	}

	compShortestDist := func(p1, p2 *path) int {
		if p1.dist < p2.dist || (p1.dist == p2.dist && len(p1.routes) < len(p2.routes)) {
			return 1
		} else if p1.dist == p2.dist && len(p1.routes) == len(p2.routes) {
			return 0
		}
		return -1
	}

	// Find paths in parallel to speed things up.
	var done sync.WaitGroup
	goAndSignal := func(comp pathComparer, bestPaths map[*city]*path) {
		done.Add(1)
		go func() {
			c.findBestPaths(comp, bestPaths)
			done.Done()
		}()
	}
	goAndSignal(compFewestHops, c.fewestHops)
	goAndSignal(compShortestDist, c.shortestDist)
	done.Wait()
}

func (c *city) findBestPaths(comp pathComparer, bestPaths map[*city]*path) {
	c.recurseBestPaths(newPath(c), make(map[*city]bool), comp, bestPaths)
}

// Completes "best" paths that begin with a given sub-path. The concept of
// "best" is determined by a caller-supplied comparison function.
//
// The function comp returns a positive number if and only if its first path
// argument is "better" than its second path argument, and it returns zero if
// and only if the first path is "equally as good" as the second path.
//
func (c *city) recurseBestPaths(p *path, visited map[*city]bool, comp pathComparer, bestPaths map[*city]*path) {

	// TODO: test

	// Compare this path to the as yet best path. (If there's no such best path
	// then this path is better by default.) If this path is better then replace
	// the best path with this path. Else, backtrack.
	curCity := p.cities[len(p.cities)-1]
	compResult := 1
	if bestPaths[curCity] != nil {
		compResult = comp(p, bestPaths[curCity])
	}
	if compResult > 0 {
		bestPaths[curCity] = copyPath(p) // replace with path
	} else {
		return // backtrack
	}

	// Traverse all routes that leave the current city, and for each route,
	// recurse.
	for adjCity, routes := range curCity.routes {
		if !visited[adjCity] {
			for _, r := range routes {
				p.appendHop(adjCity, r)
				visited[curCity] = true
				c.recurseBestPaths(p, visited, comp, bestPaths)
				visited[curCity] = false
				p.chopHop()
			}
		}
	}
}

type path struct {
	cities []*city
	routes []*route
	dist   int
}

func newPath(origCity *city) *path {
	return &path{
		cities: []*city{origCity},
	}
}

func copyPath(src *path) *path {
	return &path{
		cities: append([]*city{}, src.cities...),
		routes: append([]*route{}, src.routes...),
		dist:   src.dist,
	}
}

func (p *path) appendHop(c *city, r *route) {
	p.cities = append(p.cities, c)
	p.routes = append(p.routes, r)
	p.dist += r.dist
}

func (p *path) chopHop() {
	p.dist -= p.routes[len(p.routes)-1].dist
	p.cities = p.cities[:len(p.cities)-1]
	p.routes = p.routes[:len(p.routes)-1]
}

func (p *path) equals(other *path) bool {
	if len(p.cities) != len(other.cities) {
		return false
	}
	for i := 0; i < len(p.cities); i++ {
		if p.cities[i] != other.cities[i] {
			return false
		}
	}
	for i := 0; i < len(p.routes); i++ {
		if !p.routes[i].equals(other.routes[i]) {
			return false
		}
	}
	return true
}

func (p *path) String() (s string) {
	s = fmt.Sprintf("%q", p.cities[0].name)
	if len(p.cities)-1 != len(p.routes) {
		panic("invalid path")
	}
	for i, c := range p.cities[1:] {
		r := p.routes[i]
		s = fmt.Sprintf("%s–%d,%s–%q", s, r.dist, r.color, c.name)
	}
	return
}

type route struct {
	dist  int
	color string
}

func newRoute(dist int, color string) *route {
	return &route{
		dist:  dist,
		color: color,
	}
}

func (r *route) equals(other *route) bool {
	return r.dist == other.dist && r.color == other.color
}

type univ struct {
	cityByName map[string]*city
}

func newUniv(ents []routeEnt) (u *univ) {
	// TODO: test
	u = new(univ)
	u.cityByName = newCityMapFromRouteEntries(ents)
	for _, c := range u.cityByName {
		c.populatePaths()
	}
	return
}

func (u *univ) allCitiesAlphabetical() (cities []*city) {
	var names []string
	for n := range u.cityByName {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		cities = append(cities, u.cityByName[n])
	}
	return
}
