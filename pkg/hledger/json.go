package hledger

import (
	"encoding/json"
	"fmt"
)

// Customizes the JSON unmarshalling for JSONTag
func (t *JSONTag) UnmarshalJSON(data []byte) error {
	var raw []string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("invalid tag: %s", raw)
	}
	t.Name = raw[0]
	t.Value = raw[1]
	return nil
}
