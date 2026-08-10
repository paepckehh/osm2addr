package osm2addr

// uniform a tagSet and report whether it must be dropped (e.g. non-conform
// postcode for the target country). A true return means the entry is
// non-conform and must not be processed or exported.
func (t *tagSet) uniform() bool {
	switch t.Country {
	case "DE":
		return t.uniformDE()
	}
	return false
}
