package luarocks

import "slices"

type (
	Version struct {
		Name string
		Arch []string
	}

)

func (v *Version) String() string {
	return v.Name
}

func (v *Version) HasArch(a string) (found bool) {
	_, found = slices.BinarySearch(v.Arch, a)
	return
}

func (v *Version) AddArch(a string) {
	v.Arch = append(v.Arch, a)
	slices.Sort(v.Arch)
}

func VersionCmpFunc(v1, v2 *Version) int {
	s1, s2 := v1.Name, v2.Name
	if s1 == s2 {
		return 0
	}

	if s1 < s2 {
		return -1
	}

	return 1
}