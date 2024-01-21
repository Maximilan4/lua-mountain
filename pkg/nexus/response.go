package nexus

import (
	"encoding/json"
	"io"
)

func convertResponse[T any](r io.Reader) (*T, error) {
	var (
		model T
		d     = json.NewDecoder(r)
	)

	err := d.Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
