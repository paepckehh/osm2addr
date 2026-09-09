package osm2addr

import (
	"testing"
)

func TestUniformDEPostcode(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
			if got := string(tryNormStreetDE(street(c.in))); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestTryNormCityDE(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := string(tryNormCityDE(city(c.in))); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestCamelCaseCityDE(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
	}{
		{"frankfurt am main", "Frankfurt am Main"},
		{"garmisch-partenkirchen", "Garmisch-Partenkirchen"},
		{"berlin mitte", "Berlin Mitte"},
		{"x ot y", "X OT y"},
		{"x ii", "X II"},
		{"x ofr. y", "X OFr. y"},
		{"bad homburg v.d. höhe", "Bad Homburg V.D. Höhe"},
		{"berlin a", "Berlin a"},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			if got := camelCaseCityDE(c.in); got != c.want {
				t.Errorf("camelCaseCityDE(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestTryNormCityDETypo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
	}{
		{"X von der Y", "X vor der Y"},
		{"München in Isartal", "München im Isartal"},
		{"X in Allgäu", "X im Allgäu"},
		{"Halle / Saale", "Halle /Saale"},
		{"Untouched", "Untouched"},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			if got := tryNormCityDETypo(c.in); got != c.want {
				t.Errorf("tryNormCityDETypo(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestUniformDEStreetNormalizationOrder(t *testing.T) {
	t.Parallel()
	ts := &tagSet{Country: "DE", Postcode: "10115", City: "Berlin", Street: "Hauptstrasse"}
	if ts.uniform() {
		t.Fatal("conform entry should not be dropped")
	}
	if ts.Street != "Hauptstraße" {
		t.Errorf("street: got %q, want %q", ts.Street, "Hauptstraße")
	}
}

func TestUniformUnknownCountry(t *testing.T) {
	t.Parallel()
	ts := &tagSet{Country: "XX", Postcode: "12345", City: "Foo", Street: "Bar"}
	if ts.uniform() {
		t.Error("unknown country should not be dropped by uniform()")
	}
}

func TestIDDeterministic(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	pid := id("DE12345BerlinFoo")
	h := pid.hex()
	if len(h) != 24 {
		t.Errorf("hex place ID length: got %d, want 24", len(h))
	}
}
