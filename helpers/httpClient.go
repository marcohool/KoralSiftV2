package helpers

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"time"
)

func FetchData(endpoint string, result interface{}) error {
	client := resty.New().SetTimeout(10 * time.Second)

	log.Debug().Str("endpoint", endpoint).Msg("Fetching data")

	resp, err := client.R().SetResult(result).Get(endpoint)
	if err != nil {
		return err
	}

	if resp.IsError() {
		bodyString := resp.String()

		log.Error().
			Str("endpoint", endpoint).
			Str("status", resp.Status()).
			Str("body", bodyString).
			Msg("Failed to fetch data")

		return fmt.Errorf("API request failed with status: %s", resp.Status())
	}

	return nil
}
