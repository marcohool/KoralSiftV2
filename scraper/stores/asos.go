package stores

import (
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
	"fmt"
	"github.com/rs/zerolog/log"
)

var MensCategoryIds = []int{
	6993,
	7616,
	5668,
	3606,
	3606,
	7617,
	4208,
	3602,
	4910,
	18797,
	26090,
	14273,
	14274,
	4616,
	7078,
	7078,
	5678,
	20753,
	26776,
}

func ScrapeAsos() {
	log.Info().Msg("Starting ASOS Scraper")

	var ukMenUKProducts []models.ClothingItem

	for _, categoryId := range MensCategoryIds {
		ukMenUKProducts = append(ukMenUKProducts, GetCategoryProducts(categoryId, "Male", "GB", "GBP")...)
	}

	for _, product := range ukMenUKProducts {
		log.Info().Interface("product", product).Msg("Scraped product")
	}

	log.Info().Int("total_products", len(ukMenUKProducts)).Msg("Total products scraped")
}

func GetCategoryProducts(categoryId int, gender string, country string, currencyCode string) []models.ClothingItem {
	var data models.AsosResponse
	var clothingItems []models.ClothingItem

	var offset = 0

	for {
		err := helpers.FetchData(FormatProductsEndpoint(categoryId, offset, country, currencyCode), &data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch data")
			return nil
		}

		if len(data.Products) == 0 {
			break
		}

		for _, product := range data.Products {
			clothingItems = append(clothingItems, models.ClothingItem{
				Name:         product.Name,
				Brand:        "ASOS",
				Colours:      []string{product.Colour},
				Price:        product.Price.Current.Value,
				CurrencyCode: currencyCode,
				Gender:       gender,
				ImageUrl:     "https://" + product.ImageURL,
				SourceUrl:    "https://www.asos.com/" + product.URL,
				SourceRegion: country,
			})
		}

		offset += 200
	}

	return clothingItems
}

func FormatProductsEndpoint(categoryId int, offset int, country string, currencyCode string) string {
	return fmt.Sprintf(
		"https://www.asos.com/api/product/search/v2/categories/%d?offset=%d&includeNonPurchasableTypes=restocking&store=COM&lang=en-GB&currency=%s&channel=desktop-web&country=%s&limit=%d&excludeFacets=true",
		categoryId,
		offset,
		currencyCode,
		country,
		200,
	)
}

/*
CleansData does the following
*/
func CleanData() {
	log.Info().Msg("Cleaning ASOS data")

}
