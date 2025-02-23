package helpers

import (
	"KoralSiftV2/models"
	"fmt"
)

func MergeDuplicateClothingItems(clothingItems []models.ClothingItem) []models.ClothingItem {
	uniqueProducts := make(map[string]models.ClothingItem)

	for _, product := range clothingItems {
		if existingProduct, exists := uniqueProducts[generateProductKey(&product)]; exists {

			mergedVariants := MergeProductVariants(
				existingProduct.Variants,
				product.Variants,
			)

			existingProduct.Variants = mergedVariants
			uniqueProducts[generateProductKey(&product)] = existingProduct

			continue
		}

		uniqueProducts[generateProductKey(&product)] = product
	}

	cleanedProducts := make([]models.ClothingItem, 0, len(uniqueProducts))
	for _, product := range uniqueProducts {
		cleanedProducts = append(cleanedProducts, product)
	}

	return cleanedProducts
}

func MergeProductVariants(slice1, slice2 []models.ProductVariant) []models.ProductVariant {
	merged := make([]models.ProductVariant, 0)
	seen := make(map[string]models.ProductVariant)

	for _, variant := range append(slice1, slice2...) {
		if existing, exists := seen[variant.ImageURL]; !exists {
			merged = append(merged, variant)
			seen[variant.ImageURL] = existing
		}
	}

	return merged
}

func generateProductKey(product *models.ClothingItem) string {
	return fmt.Sprintf("%s|%s", product.Name, product.Store)
}
