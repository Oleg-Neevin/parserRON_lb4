package parserJSONWithLib

import (
	"encoding/json"
	"fmt"

	s "gitlab.se.ifmo.ru/s503298/inf_lab_4/pkg"
)

func ParseJSONWithLib(data []byte) (*s.Schedule, error) {
	var schedule s.Schedule
	err := json.Unmarshal(data, &schedule)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w", err)
	}
	return &schedule, nil
}
