package models

import "KoralSiftV2/models/enums"

type ClothingItem struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Metadata     string             `json:"metadata"`
	Brand        enums.Brand        `json:"brand"`
	Category     string             `json:"category"`
	Colours      []Colour           `json:"colours"`
	Price        float64            `json:"price"`
	CurrencyCode enums.CurrencyCode `json:"currencyCode"`
	Gender       enums.Gender       `json:"gender"`
	ImageUrl     string             `json:"imageUrl"`
	SourceUrl    string             `json:"sourceUrl"`
	SourceRegion enums.SourceRegion `json:"sourceRegion"`
}

type Colour struct {
	Name      string `json:"name"`
	Hex       string `json:"hex"`
	ImageUrl  string `json:"imageUrl"`
	SourceUrl string `json:"sourceUrl"`
}

func NewColour(name, hex, imageUrl, sourceUrl string) Colour {
	return Colour{
		Name:      name,
		Hex:       hex,
		ImageUrl:  imageUrl,
		SourceUrl: sourceUrl,
	}
}

func NewClothingItem(
	name string,
	metadata string,
	brand enums.Brand,
	colours []Colour,
	price float64,
	currencyCode enums.CurrencyCode,
	gender enums.Gender,
	imageUrl string,
	sourceUrl string,
	sourceRegion enums.SourceRegion) ClothingItem {
	return ClothingItem{
		Name:         name,
		Metadata:     metadata,
		Brand:        brand,
		Colours:      colours,
		Price:        price,
		CurrencyCode: currencyCode,
		Gender:       gender,
		ImageUrl:     imageUrl,
		SourceUrl:    sourceUrl,
		SourceRegion: sourceRegion,
	}
}
