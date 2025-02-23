package models

import "KoralSiftV2/models/enums"

type ClothingItem struct {
	Name         string             `json:"name"`
	Metadata     map[string]string  `json:"metadata"`
	Store        enums.Store        `json:"store"`
	Variants     []ProductVariant   `json:"variants"`
	CurrencyCode enums.CurrencyCode `json:"currencyCode"`
	Gender       enums.Gender       `json:"gender"`
	SourceRegion enums.SourceRegion `json:"sourceRegion"`
}

type ProductVariant struct {
	ColourName string  `json:"color"`
	ColourHex  string  `json:"hex"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
	SourceURL  string  `json:"source_url"`
}

func NewClothingItem(
	name string,
	metadata map[string]string,
	store enums.Store,
	variants []ProductVariant,
	currencyCode enums.CurrencyCode,
	gender enums.Gender,
	sourceRegion enums.SourceRegion,
) ClothingItem {
	return ClothingItem{
		Name:         name,
		Metadata:     metadata,
		Store:        store,
		Variants:     variants,
		CurrencyCode: currencyCode,
		Gender:       gender,
		SourceRegion: sourceRegion,
	}
}

func NewProductVariant(
	colourName string,
	colourHex string,
	price float64,
	imageURL string,
	sourceURL string) ProductVariant {
	return ProductVariant{
		ColourName: colourName,
		ColourHex:  colourHex,
		Price:      price,
		ImageURL:   imageURL,
		SourceURL:  sourceURL,
	}
}
