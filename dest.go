package main

import (
	"fmt"
	"math/rand"
)

type dest struct {
	city1 *city
	city2 *city
	value int
}

func countCityInDests(dests []*dest, c *city) (n int) {
	for _, d := range dests {
		if d.city1 == c || d.city2 == c {
			n++
		}
	}
	return
}

func isDestUnique(all []*dest, d *dest) bool {
	for _, x := range all {
		if x.equals(d) {
			return false
		}
	}
	return true
}

func makeDestsEqualLikely(u *univ) (dests []*dest) {

	// TODO: test

	// create slice of all cities:
	var allCities []*city
	for _, c := range u.cityByName {
		allCities = append(allCities, c)
	}

	// choose cities with equal-likely randomness:
	for len(dests) < 30 {

		index := rand.Intn(len(allCities))
		c1 := allCities[index]
		// pick city #2 such that it's at least two hops away from city #1
		var c2 *city
		for {
			index = rand.Intn(len(allCities))
			c2 = allCities[index]
			if len(c1.fewestHops[c2].routes) >= 2 {
				break
			}
		}

		d := newDest(c1, c2, c1.fewestHops[c2].dist)

		// ensure that destination is unique:
		if !isDestUnique(dests, d) {
			continue
		}

		// ensure that both cities will have at least as many unique routes as
		// destinations:
		if countCityInDests(dests, d.city1) >= len(d.city1.routes) {
			continue
		}
		if countCityInDests(dests, d.city2) >= len(d.city2.routes) {
			continue
		}

		// destination is OK:
		dests = append(dests, newDest(c1, c2, c1.fewestHops[c2].dist))
	}

	return
}

func newDest(c1, c2 *city, value int) *dest {
	return &dest{
		city1: c1,
		city2: c2,
		value: value,
	}
}

func newDestsFromDestEntries(u *univ, ents []destEnt) (s []*dest, err error) {
	// TODO: test
	for _, ent := range ents {
		c1 := u.cityByName[ent.name1]
		if c1 == nil {
			return nil, fmt.Errorf("error creating destination: city %q doesn't exist", ent.name1)
		}
		c2 := u.cityByName[ent.name2]
		if c2 == nil {
			return nil, fmt.Errorf("error creating destination: city %q doesn't exist", ent.name2)
		}
		s = append(s, newDest(c1, c2, ent.value))
	}
	return
}

func (d *dest) equals(other *dest) bool {
	return (d.city1 == other.city1 && d.city2 == other.city2) || (d.city1 == other.city2 && d.city2 == other.city1)
}
