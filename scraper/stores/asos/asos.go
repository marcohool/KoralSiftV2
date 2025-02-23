package asos

import (
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
	"KoralSiftV2/models/enums"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
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
		ukMenUKProducts = append(
			ukMenUKProducts,
			GetCategoryProducts(categoryId, enums.Male, enums.UK, enums.GBP)...)
	}

	log.Info().Int("total_products", len(ukMenUKProducts)).Msg("Total products scraped")

	err := helpers.SaveSliceToJSONFile(ukMenUKProducts, "asos-dirty")
	if err != nil {
		log.Error().Err(err).Msg("Failed to save ASOS data")
	}

	cleanedClothingItems := CleanASOSData(ukMenUKProducts)

	log.Info().
		Int("total_cleaned_products", len(cleanedClothingItems)).
		Msg("Total cleaned products")

	err = helpers.SaveSliceToJSONFile(cleanedClothingItems, "asos")
	if err != nil {
		log.Error().Err(err).Msg("Failed to save ASOS data")
	}

}

func GetCategoryProducts(
	categoryId int,
	gender enums.Gender,
	country enums.SourceRegion,
	currencyCode enums.CurrencyCode,
) []models.ClothingItem {
	var data models.AsosResponse
	var clothingItems []models.ClothingItem

	var offset = 0
	var limit = 200

	var countryMap = map[enums.SourceRegion]string{
		enums.UK: "GB",
	}

	for {
		var productsEndpoint = fmt.Sprintf(
			"https://www.asos.com/api/product/search/v2/categories/%d?offset=%d&includeNonPurchasableTypes=restocking&store=COM&lang=en-GB&currency=%s&channel=desktop-web&country=%s&limit=%d&excludeFacets=true",
			categoryId,
			offset,
			currencyCode,
			countryMap[country],
			limit,
		)

		err := helpers.FetchData(productsEndpoint, &data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch data")
			return nil
		}

		if len(data.Products) == 0 {
			break
		}

		for _, product := range data.Products {
			var sourceUrl = fmt.Sprintf("https://www.asos.com/%s", product.URL)
			var imageUrl = fmt.Sprintf("https://%s", product.ImageURL)

			clothingItems = append(clothingItems, models.NewClothingItem(
				product.Name,
				nil,
				enums.ASOS,
				[]models.ProductVariant{
					models.NewProductVariant(
						product.Colour,
						"",
						product.Price.Current.Value,
						imageUrl,
						sourceUrl,
					),
				},
				currencyCode,
				gender,
				country,
			))
		}

		offset += limit
	}

	return clothingItems
}

func CleanASOSData(clothingItems []models.ClothingItem) []models.ClothingItem {
	log.Info().Msg("Cleaning ASOS data")

	for i := range clothingItems {
		clothingItems[i].Name = removeColourFromName(clothingItems[i].Name)
	}

	uniqueProducts := helpers.MergeDuplicateClothingItems(clothingItems)

	return uniqueProducts
}

func removeColourFromName(name string) string {
	parts := strings.Split(strings.ToLower(name), " in ")
	if len(parts) == 2 {
		return parts[0]
	}
	return name
}
