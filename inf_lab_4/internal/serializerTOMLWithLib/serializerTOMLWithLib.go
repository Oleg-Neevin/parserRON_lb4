package serializerTOMLWithLib

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"

	s "gitlab.se.ifmo.ru/s503298/inf_lab_4/pkg"
)

func SerializeTOMLWithLib(schedule *s.Schedule) (string, error) {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)

	err := encoder.Encode(schedule)
	if err != nil {
		return "", fmt.Errorf("TOML encode error: %w", err)
	}

	return buf.String(), nil
}
