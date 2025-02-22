package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func SaveSliceToJSONFile(data interface{}, filename string) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05") // YYYY-MM-DD_HH-MM-SS format
	formattedFilename := fmt.Sprintf("data/scans/%s_%s.json", filename, timestamp)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	err = os.WriteFile(formattedFilename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	log.Info().Str("filename", formattedFilename).Msg("Saved JSON data to file")

	return nil
}
