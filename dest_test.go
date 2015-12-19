package main

import (
	"testing"
)

func TestDestEquals(t *testing.T) {

	check := func(d1, d2 *dest, exp bool) {
		if d1.equals(d2) != exp {
			t.Errorf("with %v and %v, expected %v but got %v", d1, d2, exp, !exp)
		}
		if d2.equals(d1) != exp {
			t.Errorf("with %v and %v, expected %v but got %v", d2, d1, exp, !exp)
		}
	}

	c1 := newCity("alpha")
	c2 := newCity("bravo")
	c3 := newCity("charlie")

	// destinations are same:
	check(newDest(c1, c2, 1), newDest(c1, c2, 1), true)
	check(newDest(c1, c2, 1), newDest(c1, c2, 2), true)
	check(newDest(c1, c2, 1), newDest(c2, c1, 1), true)
	check(newDest(c1, c2, 1), newDest(c2, c1, 2), true)

	// destinations are different:
	check(newDest(c1, c3, 1), newDest(c1, c2, 1), false)
	check(newDest(c1, c3, 1), newDest(c1, c2, 2), false)
	check(newDest(c1, c3, 1), newDest(c2, c1, 1), false)
	check(newDest(c1, c3, 1), newDest(c2, c1, 2), false)
}

func TestIsDestUnique(t *testing.T) {
	c1 := newCity("alpha")
	c2 := newCity("bravo")
	c3 := newCity("charlie")
	dests := []*dest{
		newDest(c1, c2, 1),
		newDest(c1, c3, 1),
	}
	if isDestUnique(dests, newDest(c1, c2, 2)) {
		t.Errorf("expected non-unique but got unique")
	}
	if !isDestUnique(dests, newDest(c2, c3, 3)) {
		t.Errorf("expected unique but got non-unique")
	}
}

func TestCountCityInDests(t *testing.T) {

	check := func(dests []*dest, c *city, exp int) {
		got := countCityInDests(dests, c)
		if exp != got {
			t.Errorf("with %v, %v, expected %v but got %v", dests, c, exp, got)
		}
	}

	c1 := newCity("alpha")
	c2 := newCity("bravo")
	c3 := newCity("charlie")
	c4 := newCity("delta")
	dests := []*dest{
		newDest(c1, c2, 1),
		newDest(c1, c3, 1),
	}
	check(dests, c1, 2)
	check(dests, c2, 1)
	check(dests, c3, 1)
	check(dests, c4, 0)
}
