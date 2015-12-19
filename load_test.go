package main

import (
	"strings"
	"testing"
)

func TestLoadDestEntriesParser(t *testing.T) {
	type tc struct {
		inText  string
		expEnts []destEnt
	}

	tcs := []tc{
		// no whitespace:
		{"alpha-bravo:2\nalpha-charlie:2\n", []destEnt{
			destEnt{"alpha", "bravo", 2},
			destEnt{"alpha", "charlie", 2},
		}},
		// "normal" whitespace:
		{"alpha - bravo: 2\nalpha - charlie: 2\n", []destEnt{
			destEnt{"alpha", "bravo", 2},
			destEnt{"alpha", "charlie", 2},
		}},
		// whitespace before first city
		{" \t alpha - bravo: 2\n \t alpha - charlie: 2\n", []destEnt{
			destEnt{"alpha", "bravo", 2},
			destEnt{"alpha", "charlie", 2},
		}},

		// tabs instead of spaces:
		{"alpha\t-\tbravo:\t2\nalpha\t-\tcharlie:\t2\n", []destEnt{
			destEnt{"alpha", "bravo", 2},
			destEnt{"alpha", "charlie", 2},
		}},
		// empty lines:
		{"alpha - bravo: 2\n\n\nalpha - charlie: 2\n\n\n", []destEnt{
			destEnt{"alpha", "bravo", 2},
			destEnt{"alpha", "charlie", 2},
		}},
		// no end-of-line on last line:
		{"alpha - bravo: 2\nalpha - charlie: 2", []destEnt{
			destEnt{"alpha", "bravo", 2},
			destEnt{"alpha", "charlie", 2},
		}},
		// empty input:
		{"", []destEnt{}},
	}

	// run test cases:
	for _, tc := range tcs {
		if dests, err := loadDestEntries(strings.NewReader(tc.inText)); err != nil {
			t.Errorf("got error loading %q: %s", tc.inText, err)
		} else if len(dests) != len(tc.expEnts) {
			t.Errorf("expected %v destination(s) but got %v (%q, %v)", len(tc.expEnts), len(dests), tc.inText, dests)
		} else {
			for i, exp := range tc.expEnts {
				got := dests[i]
				if exp != got {
					t.Errorf("expected destination %v to be %v but got %v", i, exp, got)
				}
			}
		}
	}
}

func TestLoadDestEntriesReal(t *testing.T) {
	mustLoadDestEntriesFromFile("destinations.dat")
}

func TestLoadRouteEntriesParser(t *testing.T) {
	type tc struct {
		inText  string
		expEnts []routeEnt
	}

	tcs := []tc{
		// no whitespace:
		{"alpha-bravo:2 blue,2 orange\nalpha-charlie:2 wild\n", []routeEnt{
			routeEnt{"alpha", "bravo", 2, "blue"},
			routeEnt{"alpha", "bravo", 2, "orange"},
			routeEnt{"alpha", "charlie", 2, "wild"},
		}},
		// "normal" whitespace:
		{"alpha - bravo: 2 blue, 2 orange\nalpha - charlie: 2 wild\n", []routeEnt{
			routeEnt{"alpha", "bravo", 2, "blue"},
			routeEnt{"alpha", "bravo", 2, "orange"},
			routeEnt{"alpha", "charlie", 2, "wild"},
		}},
		// whitespace before first city
		{" \t alpha - bravo: 2 blue, 2 orange\n \t alpha - charlie: 2 wild\n", []routeEnt{
			routeEnt{"alpha", "bravo", 2, "blue"},
			routeEnt{"alpha", "bravo", 2, "orange"},
			routeEnt{"alpha", "charlie", 2, "wild"},
		}},
		// tabs instead of spaces:
		{"alpha\t-\tbravo:\t2\tblue,\t2\torange\nalpha\t-\tcharlie:\t2\twild\n", []routeEnt{
			routeEnt{"alpha", "bravo", 2, "blue"},
			routeEnt{"alpha", "bravo", 2, "orange"},
			routeEnt{"alpha", "charlie", 2, "wild"},
		}},
		// empty lines:
		{"\n\nalpha - bravo: 2 blue, 2 orange\n\n\nalpha - charlie: 2 wild\n\n\n", []routeEnt{
			routeEnt{"alpha", "bravo", 2, "blue"},
			routeEnt{"alpha", "bravo", 2, "orange"},
			routeEnt{"alpha", "charlie", 2, "wild"},
		}},
		// no end-of-line on last line:
		{"alpha - bravo: 2 wild\nalpha - charlie: 2 wild", []routeEnt{
			routeEnt{"alpha", "bravo", 2, "wild"},
			routeEnt{"alpha", "charlie", 2, "wild"},
		}},
		// empty input:
		{"", []routeEnt{}},
	}

	// run test cases:
	for _, tc := range tcs {
		if routes, err := loadRouteEntries(strings.NewReader(tc.inText)); err != nil {
			t.Errorf("got error loading %q: %s", tc.inText, err)
		} else if len(routes) != len(tc.expEnts) {
			t.Errorf("expected %v route(s) but got %v (%q, %v)", len(tc.expEnts), len(routes), tc.inText, routes)
		} else {
			for i, exp := range tc.expEnts {
				got := routes[i]
				if exp != got {
					t.Errorf("expected route %v to be %v but got %v", i, exp, got)
				}
			}
		}
	}
}

func TestLoadRouteEntriesReal(t *testing.T) {
	mustLoadRouteEntriesFromFile("routes.dat")
}
