package models

type ClothingItem struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Brand        string   `json:"brand"`
	Category     string   `json:"category"`
	Colours      []string `json:"colours"`
	Price        float64  `json:"price"`
	CurrencyCode string   `json:"currencyCode"`
	Gender       string   `json:"gender"`
	ImageUrl     string   `json:"imageUrl"`
	SourceUrl    string   `json:"sourceUrl"`
	SourceRegion string   `json:"sourceRegion"`
}
