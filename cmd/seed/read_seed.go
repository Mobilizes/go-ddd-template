package seed

import (
	"encoding/json"
	"fmt"
	"os"
)

func readSeedJSON(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	return nil
}
