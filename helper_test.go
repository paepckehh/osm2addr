package osm2addr

import (
	"testing"
)

func TestIsASCII(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"Berlin", true},
		{"K\xf6ln", false},
		{"Stra\xdfe", false},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			if got := isASCII(c.in); got != c.want {
				t.Errorf("isASCII(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestIsLatin1(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"Berlin", true},
		{"Köln", true},
		{"Straße", true},
		{"Москва", false},
		{"日本", false},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			if got := isLatin1(c.in); got != c.want {
				t.Errorf("isLatin1(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestMakeCapitalLetter(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"lowercase", "berlin", "Berlin"},
		{"already capital", "Berlin", "Berlin"},
		{"uppercase umlaut", "über", "Über"},
		{"empty", "", ""},
		{"single letter", "a", "A"},
		{"valid replacement rune is not an error", "\uFFFD", "\uFFFD"},
		{"invalid utf8 is returned unchanged", "\xff\xfe", "\xff\xfe"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if got := makeCapitalLetter(c.in); got != c.want {
				t.Errorf("makeCapitalLetter(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestCamelCaseSep(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, sep, want string
	}{
		{"bad-homburg", "-", "bad-Homburg"},
		{"bad-homburg", "(", "bad-homburg"},
		{"unter den linden (west)", "(", "unter den linden (West)"},
		{"a.b.c", ".", "a.B.C"},
		{"no-separator", "-", "no-Separator"},
	}
	for _, c := range cases {
		t.Run(c.in+"/"+c.sep, func(t *testing.T) {
			t.Parallel()
			if got := camelCaseSep(c.in, c.sep); got != c.want {
				t.Errorf("camelCaseSep(%q, %q) = %q, want %q", c.in, c.sep, got, c.want)
			}
		})
	}
}

func TestCamelCaseSeps(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
	}{
		{"garmisch-partenkirchen", "garmisch-Partenkirchen"},
		{"frankfurt/main", "frankfurt/Main"},
		{"bad homburg v.d.höhe", "bad homburg v.D.Höhe"},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			if got := camelCaseSeps(c.in); got != c.want {
				t.Errorf("camelCaseSeps(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestHu(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   int
		want string
	}{
		{0, "          0"},
		{42, "         42"},
		{278381, "    278.381"},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			t.Parallel()
			got := hu(c.in)
			if len(got) != 11 {
				t.Errorf("hu(%d) length = %d, want 11 (%q)", c.in, len(got), got)
			}
			if got != c.want {
				t.Errorf("hu(%d) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
