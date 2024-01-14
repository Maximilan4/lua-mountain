package luarocks

import (
	"cmp"
	"golang.org/x/exp/slices"
)

type (
	Rock struct {
		Name string
		Versions []*Version
	}
	
	RocksList []*Rock
)
func (rl *RocksList) Add(rockName, version, arch string) {
	rock := rl.Search(rockName)
	crl := *rl
	if rock == nil {
		rock = &Rock{Name: rockName, Versions: make([]*Version, 0, 3)}
		crl = append(crl, rock)
	}

	v := rock.SearchVersion(version)
	if v == nil {
		v = &Version{version, make([]string, 0, 3)}
		v.AddArch(arch)
		rock.AddVersion(v)
	} else if !v.HasArch(arch) {
		v.AddArch(arch)
	}

	slices.SortFunc(crl, func(a, b *Rock) int {
		return cmp.Compare(a.Name, b.Name)
	})
	*rl = crl
}

func (rl *RocksList) Search(rockName string) *Rock {
	crl := *rl
	if len(crl) == 0 {
		return nil
	}

	pos, found := slices.BinarySearchFunc(crl, rockName, func(rock *Rock, s string) int {
		return cmp.Compare(rock.Name, s)
	})

	if !found {
		return nil
	}

	return crl[pos]
}

func NewRock(name string) *Rock {
	return &Rock{Name: name}
}

func (p *Rock) String() string {
	return p.Name
}

func (p *Rock) SearchVersion(v string) *Version {
	if len(p.Versions) == 0 {
		return nil
	}

	pos, found := slices.BinarySearchFunc(p.Versions, v, func(version *Version, s string) int {
		return cmp.Compare(version.Name, s)
	})
	if !found {
		return nil
	}

	return p.Versions[pos]
}

func (p *Rock) AddVersion(v *Version) {
	p.Versions = append(p.Versions, v)
	slices.SortFunc(p.Versions, VersionCmpFunc)
}
