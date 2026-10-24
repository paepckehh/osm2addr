package decoder

import "testing"

func TestStrAtBounds(t *testing.T) {
	t.Parallel()
	table := []string{"zero", "one"}
	cases := []struct {
		i    int
		want string
	}{
		{-1, ""},
		{0, "zero"},
		{1, "one"},
		{2, ""},
	}
	for _, c := range cases {
		if got := strAt(table, c.i); got != c.want {
			t.Errorf("strAt(%d) = %q, want %q", c.i, got, c.want)
		}
	}
}

func TestDecodeTagsOutOfRange(t *testing.T) {
	t.Parallel()
	c := &blockContext{strings: []string{"", "key", "value"}}

	tags := c.decodeTags([]uint32{1, 99}, []uint32{2, 99})
	if tags["key"] != "value" {
		t.Errorf("in-range tag lost: %v", tags)
	}

	tic := c.newTagsContext([]int32{1, 2, 98, 99, 0})
	got := tic.decodeTags()
	if got["key"] != "value" {
		t.Errorf("in-range dense tag lost: %v", got)
	}
}
