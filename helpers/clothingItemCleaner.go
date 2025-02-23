package helpers

//func MergeDuplicatesBySourceUrl(clothingItems []models.ClothingItem) []models.ClothingItem {
//	uniqueProducts := make(map[string]models.ClothingItem)
//
//	for _, product := range clothingItems {
//		if _, exists := uniqueProducts[product.SourceUrl]; exists {
//			mergedColours := MergeColours(
//				uniqueProducts[product.SourceUrl].Colours,
//				product.Colours,
//			)
//
//			product.Colours = mergedColours
//			uniqueProducts[product.SourceUrl] = product
//
//			continue
//		}
//		uniqueProducts[product.SourceUrl] = product
//	}
//
//	cleanedProducts := make([]models.ClothingItem, 0, len(uniqueProducts))
//	for _, product := range uniqueProducts {
//		cleanedProducts = append(cleanedProducts, product)
//	}
//
//	return cleanedProducts
//}
