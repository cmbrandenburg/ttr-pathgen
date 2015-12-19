package main

import (
	"fmt"
	"testing"
)

// Returns an incomplete universe containing all the cities and routes but
// without any of the calculated paths.
func newUnivNoPaths(ents []routeEnt) (u *univ) {
	u = new(univ)
	u.cityByName = newCityMapFromRouteEntries(ents)
	return
}

func TestNewCityMapFromRouteEntries(t *testing.T) {
	type tc struct {
		ents         []routeEnt
		expCityNames []string
		expRoutes    []string
	}

	tcs := []tc{

		// case: one route
		tc{
			[]routeEnt{
				routeEnt{name1: "alpha", name2: "bravo", dist: 3, color: "wild"},
			},
			[]string{
				"alpha",
				"bravo",
			},
			[]string{
				"alpha bravo 3 wild",
				"bravo alpha 3 wild",
			},
		},

		// case: double route between two cities
		tc{
			[]routeEnt{
				routeEnt{name1: "alpha", name2: "bravo", dist: 3, color: "blue"},
				routeEnt{name1: "alpha", name2: "bravo", dist: 3, color: "orange"},
			},
			[]string{
				"alpha",
				"bravo",
			},
			[]string{
				"alpha bravo 3 blue",
				"alpha bravo 3 orange",
				"bravo alpha 3 blue",
				"bravo alpha 3 orange",
			},
		},

		// case: routes from one city to many cities
		tc{
			[]routeEnt{
				routeEnt{name1: "alpha", name2: "bravo", dist: 3, color: "blue"},
				routeEnt{name1: "alpha", name2: "bravo", dist: 3, color: "orange"},
				routeEnt{name1: "alpha", name2: "charlie", dist: 2, color: "wild"},
			},
			[]string{
				"alpha",
				"bravo",
				"charlie",
			},
			[]string{
				"alpha bravo 3 blue",
				"alpha bravo 3 orange",
				"alpha charlie 2 wild",
				"bravo alpha 3 blue",
				"bravo alpha 3 orange",
				"charlie alpha 2 wild",
			},
		},

		// case: completely connected graph of three cities
		tc{
			[]routeEnt{
				routeEnt{name1: "alpha", name2: "bravo", dist: 4, color: "wild"},
				routeEnt{name1: "alpha", name2: "charlie", dist: 6, color: "wild"},
				routeEnt{name1: "bravo", name2: "charlie", dist: 1, color: "wild"},
			},
			[]string{
				"alpha",
				"bravo",
				"charlie",
			},
			[]string{
				"alpha bravo 4 wild",
				"alpha charlie 6 wild",
				"bravo alpha 4 wild",
				"bravo charlie 1 wild",
				"charlie alpha 6 wild",
				"charlie bravo 1 wild",
			},
		},
	}

	for _, tc := range tcs {
		got := newCityMapFromRouteEntries(tc.ents)
		expNames := make(map[string]bool)
		for _, v := range tc.expCityNames {
			expNames[v] = true
		}
		expRoutes := make(map[string]bool)
		for _, v := range tc.expRoutes {
			expRoutes[v] = true
		}
		// Check that all expected city names are in the map.
		for _, expName := range tc.expCityNames {
			if got[expName] == nil {
				t.Errorf("missing city %q (%v)", expName, tc)
			}
		}
		// Check that all actual city names are expected.
		for gotName := range got {
			if !expNames[gotName] {
				t.Errorf("got unexpected city %q (%v)", gotName, tc)
			}
		}
		// Check that all expected routes are in the map.
		for _, v := range tc.expRoutes {
			var c1Name string
			var c2Name string
			var dist int
			var color string
			if _, err := fmt.Sscan(v, &c1Name, &c2Name, &dist, &color); err != nil {
				t.Fatalf("scan failure: %s (%q)", err, v)
			}
			c1 := got[c1Name]
			c2 := got[c2Name]
			found := false
			for _, r := range c1.routes[c2] {
				if r.dist == dist && r.color == color {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("missing route {%q - %q, %d, %s} (%v)", c1.name, c2.name, dist, color, tc)
			}
		}
		// Check that all actual routes are expected.
		for origName, origCity := range got {
			for tgtCity, routes := range origCity.routes {
				tgtName := tgtCity.name
				for _, r := range routes {
					k := fmt.Sprintf("%s %s %d %s", origName, tgtName, r.dist, r.color)
					if !expRoutes[k] {
						t.Errorf("got unexpected route {%q - %q, %d, %s} (%v)", origName, tgtName, r.dist, r.color, tc)
					}
				}
			}
		}
	}
}

func TestPathAppendAndChop(t *testing.T) {

	c0 := newCity("alpha")
	c1 := newCity("bravo")
	c2 := newCity("bravo")
	r1 := newRoute(5, "green")
	r2 := newRoute(4, "wild")
	p := newPath(c0)

	check := func(expCities []*city, expRoutes []*route) {
		if len(p.cities) != len(expCities) {
			t.Fatalf("expected %d cities but got %d", len(expCities), len(p.cities))
		}
		for i, got := range p.cities {
			exp := expCities[i]
			if got != exp {
				t.Fatalf("at index %d, expected %q but got %q", i, exp.name, got.name)
			}
		}
		expDist := 0
		if len(p.routes) != len(expRoutes) {
			t.Fatalf("expected %d routes but got %d", len(expRoutes), len(p.routes))
		}
		for i, got := range p.routes {
			exp := expRoutes[i]
			if got != exp {
				t.Fatalf("at index %d, expected {%d, %s} but got {%d, %s}", i, exp.dist, exp.color, got.dist, got.color)
			}
			expDist += exp.dist
		}
		if p.dist != expDist {
			t.Fatalf("expected distance %d but got %d", expDist, p.dist)
		}
	}

	check([]*city{c0}, []*route{})
	p.appendHop(c1, r1)
	check([]*city{c0, c1}, []*route{r1})
	p.appendHop(c2, r2)
	check([]*city{c0, c1, c2}, []*route{r1, r2})
	p.chopHop()
	check([]*city{c0, c1}, []*route{r1})
	p.chopHop()
	check([]*city{c0}, []*route{})
}

func TestPathEquals(t *testing.T) {
	c1 := newCity("alpha")
	c2 := newCity("bravo")
	c3 := newCity("charlie")
	var p1, p2 *path

	p1 = newPath(c1)

	// check: same path object is same
	if !p1.equals(p1) {
		t.Errorf("path object reported unequal to self")
	}

	// check: one-hop paths that are the same
	p1 = newPath(c1)
	p2 = newPath(c1)
	if !p1.equals(p2) {
		t.Errorf("similar one-hop paths reported unequal")
	}

	// check: one-hop paths that are different
	p1 = newPath(c1)
	p2 = newPath(c2)
	if p1.equals(p2) {
		t.Errorf("different one-hop paths reported equal")
	}

	// check: multiple-paths that are the same
	p1 = newPath(c1)
	p1.appendHop(c2, newRoute(2, "red"))
	p2 = newPath(c1)
	p2.appendHop(c2, newRoute(2, "red"))
	if !p1.equals(p2) {
		t.Errorf("similar multi-hop paths reported unequal")
	}

	// check: multiple-paths that have different hops
	p1 = newPath(c1)
	p1.appendHop(c2, newRoute(2, "red"))
	p2 = newPath(c1)
	p2.appendHop(c3, newRoute(2, "red"))
	if p1.equals(p2) {
		t.Errorf("similar multi-hop paths reported equal")
	}

	// check: multiple-paths that have different routes
	p1 = newPath(c1)
	p1.appendHop(c2, newRoute(2, "red"))
	p2 = newPath(c1)
	p2.appendHop(c2, newRoute(2, "blue"))
	if p1.equals(p2) {
		t.Errorf("similar multi-hop paths reported equal")
	}
}

func TestCityFindBestPaths(t *testing.T) {
}

func BenchmarkNewDefaultUniverse(b *testing.B) {
	routeEnts := mustLoadRouteEntriesFromFile("routes.dat")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newUniv(routeEnts)
	}
}
