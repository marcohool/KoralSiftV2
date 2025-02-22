package models

import "KoralSiftV2/models/enums"

type ClothingItem struct {
	Name         string             `json:"name"`
	Metadata     map[string]string  `json:"metadata"`
	Store        enums.Store        `json:"store"`
	Price        float64            `json:"price"`
	Variants     []ColourVariant    `json:"variants"`
	CurrencyCode enums.CurrencyCode `json:"currencyCode"`
	Gender       enums.Gender       `json:"gender"`
	SourceRegion enums.SourceRegion `json:"sourceRegion"`
}

type ColourVariant struct {
	Color     string `json:"color"`
	ImageURL  string `json:"image_url"`
	SourceURL string `json:"source_url"`
}

func NewClothingItem(
	name string,
	metadata map[string]string,
	store enums.Store,
	price float64,
	variants []ColourVariant,
	currencyCode enums.CurrencyCode,
	gender enums.Gender,
	sourceRegion enums.SourceRegion,
) ClothingItem {
	return ClothingItem{
		Name:         name,
		Metadata:     metadata,
		Store:        store,
		Price:        price,
		Variants:     variants,
		CurrencyCode: currencyCode,
		Gender:       gender,
		SourceRegion: sourceRegion,
	}
}

func NewColourVariant(
	color string,
	imageURL string,
	sourceURL string) ColourVariant {
	return ColourVariant{
		Color:     color,
		ImageURL:  imageURL,
		SourceURL: sourceURL,
	}
}
