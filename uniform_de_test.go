package osm2addr

import (
	"testing"
)

func TestUniformDEPostcode(t *testing.T) {
	cases := []struct {
		in   string
		want string
		drop bool
	}{
		{"12345", "12345", false},
		{"1234", "01234", false},
		{"0123", "", true},
		{"123", "", true},
		{"123456", "", true},
		{"ABCDE", "", true},
		{"-1234", "", true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			ts := &tagSet{Country: "DE", Postcode: postcode(c.in), City: "Berlin", Street: "Teststr."}
			drop := ts.uniform()
			if drop != c.drop {
				t.Fatalf("drop: got %v, want %v", drop, c.drop)
			}
			if !c.drop && string(ts.Postcode) != c.want {
				t.Errorf("postcode: got %q, want %q", ts.Postcode, c.want)
			}
		})
	}
}

func TestUniformDENonLatin1(t *testing.T) {
	ts := &tagSet{Country: "DE", Postcode: "12345", City: "Москва", Street: "Foo"}
	if !ts.uniform() {
		t.Error("non-Latin1 city should be dropped")
	}
	ts = &tagSet{Country: "DE", Postcode: "12345", City: "Berlin", Street: "Москва"}
	if !ts.uniform() {
		t.Error("non-Latin1 street should be dropped")
	}
}

func TestUniformDENormalizesCityStreet(t *testing.T) {
	ts := &tagSet{Country: "DE", Postcode: "12345", City: "frankfurt a.d. oder", Street: "Hauptstrasse"}
	if ts.uniform() {
		t.Fatal("conform entry should not be dropped")
	}
	if ts.City != "Frankfurt an der Oder" {
		t.Errorf("city: got %q, want %q", ts.City, "Frankfurt an der Oder")
	}
	if ts.Street != "Hauptstraße" {
		t.Errorf("street: got %q, want %q", ts.Street, "Hauptstraße")
	}
}

func TestTryNormStreetDE(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hauptstrasse", "Hauptstraße"},
		{"Hauptstr.", "Haupt straße"},
		{"Test Strasse", "Test Straße"},
		{"Test Str.", "Test Straße"},
		{"test strasse", "test Straße"},
		{"test str.", "test Straße"},
		{"", ""},
		{"Already Straße", "Already Straße"},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if got := string(tryNormStreetDE(street(c.in))); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestTryNormCityDE(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"frankfurt a.d. oder", "Frankfurt an der Oder"},
		{"neustadt a.d. donau", "Neustadt an der Donau"},
		{"Berlin", "Berlin"},
		{"berlin", "Berlin"},
		{"Köln", "Köln"},
		{"bad homburg v.d. höhe", "Bad Homburg vor der Höhe"},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if got := string(tryNormCityDE(city(c.in))); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestUniformUnknownCountry(t *testing.T) {
	ts := &tagSet{Country: "XX", Postcode: "12345", City: "Foo", Street: "Bar"}
	if ts.uniform() {
		t.Error("unknown country should not be dropped by uniform()")
	}
}

func TestIDDeterministic(t *testing.T) {
	a := id("DE12345BerlinFoo")
	b := id("DE12345BerlinFoo")
	if a != b {
		t.Error("place ID must be deterministic for identical input")
	}
	c := id("DE12345BerlinBar")
	if a == c {
		t.Error("different input must yield different place ID")
	}
}

func TestPlaceIDHex(t *testing.T) {
	pid := id("DE12345BerlinFoo")
	h := pid.hex()
	if len(h) != 24 {
		t.Errorf("hex place ID length: got %d, want 24", len(h))
	}
}
