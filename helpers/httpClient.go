package helpers

import (
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
		log.Error().Str("endpoint", endpoint).Str("status", resp.Status()).Msg("Failed to fetch data")
	}

	return nil
}
