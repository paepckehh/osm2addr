package osm2addr

// uniform a tagSet
func (t *TagSET) uniform() int {
	switch t.Country {
	case "DE":
		return t.uniformDE()
	}
	return 0
}
