package main

import (
	"KoralSiftV2/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	setupLogger()

	log.Info().Msg("Starting application")

	scraper.RunScraper()
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				return "\033[36m" + ll + "\033[0m" // 🔹 Set DEBUG to Cyan instead of Red
			case "info":
				return "\033[32m" + ll + "\033[0m" // ✅ Green for Info
			case "warn":
				return "\033[33m" + ll + "\033[0m" // 🟡 Yellow for Warn
			case "error":
				return "\033[31m" + ll + "\033[0m" // 🔴 Red for Error
			default:
				return ll
			}
		}
		return i.(string)
	}

	log.Logger = log.Output(output)
}
