package stores

import (
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
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
		ukMenUKProducts = append(ukMenUKProducts, GetCategoryProducts(categoryId, "Male", "GB", "GBP")...)
	}

	for _, product := range ukMenUKProducts {
		log.Info().Interface("product", product).Msg("Scraped product")
	}

	log.Info().Int("total_products", len(ukMenUKProducts)).Msg("Total products scraped")

	cleanedClothingItems := CleanASOSData(ukMenUKProducts)

	log.Info().Int("total_cleaned_products", len(cleanedClothingItems)).Msg("Total cleaned products")

	for _, product := range cleanedClothingItems {
		log.Info().Interface("product", product).Msg("Cleaned product")
	}
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
				Colours:      []models.Colour{{Name: product.Colour, ImageUrl: product.ImageURL}},
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

// CleanASOSData Removes duplicates and merges color variations
func CleanASOSData(clothingItems []models.ClothingItem) []models.ClothingItem {
	log.Info().Msg("Cleaning ASOS data")

	// Step 1: Remove exact duplicates based on SourceUrl
	uniqueProductsBySourceUrl := helpers.MergeDuplicatesBySourceUrl(clothingItems)

	// Step 2: Remove color variations from product name
	uniqueProductsByName := make(map[string]models.ClothingItem)
	for _, product := range uniqueProductsBySourceUrl {
		baseName := RemoveColourFromName(product.Name)
		product.Name = baseName

		uniqueProductsByName[product.SourceUrl] = product
	}

	// Step 3: Merge products with different colors but same base name
	nameMap := make(map[string]*models.ClothingItem)
	for _, product := range uniqueProductsByName {
		if existingProduct, found := nameMap[product.Name]; found {
			existingProduct.Colours = helpers.MergeColours(existingProduct.Colours, product.Colours)
		} else {
			nameMap[product.Name] = &product
		}
	}

	// Step 4: Convert map back to slice
	cleanedProducts := make([]models.ClothingItem, 0, len(nameMap))
	for _, product := range nameMap {
		cleanedProducts = append(cleanedProducts, *product)
	}

	return cleanedProducts
}

func RemoveColourFromName(name string) string {
	parts := strings.Split(name, " in ")
	if len(parts) == 2 {
		return parts[0]
	}
	return name
}
