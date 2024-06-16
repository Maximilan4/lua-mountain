package nexus

import (
	"encoding/json"
	"io"
)

func jsonTo[T any](r io.Reader) (*T, error) {
	var (
		o T
		d = json.NewDecoder(r)
	)

	err := d.Decode(&o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}
