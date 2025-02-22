package models

type AsosResponse struct {
	SearchTerm   string    `json:"searchTerm"`
	CategoryName string    `json:"categoryName"`
	ItemCount    int       `json:"itemCount"`
	RedirectUrl  string    `json:"redirectUrl"`
	Products     []Product `json:"products"`
}

type Product struct {
	ID                  int             `json:"id"`
	Name                string          `json:"name"`
	Price               PriceDetails    `json:"price"`
	Colour              string          `json:"colour"`
	ColourWayID         int             `json:"colourWayId"`
	BrandName           string          `json:"brandName"`
	HasVariantColours   bool            `json:"hasVariantColours"`
	HasMultiplePrices   bool            `json:"hasMultiplePrices"`
	GroupID             *int            `json:"groupId"`
	ProductCode         int             `json:"productCode"`
	ProductType         string          `json:"productType"`
	URL                 string          `json:"url"`
	ImageURL            string          `json:"imageUrl"`
	AdditionalImageURLs []string        `json:"additionalImageUrls"`
	VideoURL            *string         `json:"videoUrl"`
	ShowVideo           bool            `json:"showVideo"`
	IsSellingFast       bool            `json:"isSellingFast"`
	IsRestockingSoon    bool            `json:"isRestockingSoon"`
	IsPromotion         bool            `json:"isPromotion"`
	SponsoredCampaignID *int            `json:"sponsoredCampaignId"`
	FacetGroupings      []FacetGrouping `json:"facetGroupings"`
	Advertisement       *string         `json:"advertisement"`
	EarlyAccess         *string         `json:"earlyAccess"`
}

type PriceDetails struct {
	Current                 PriceInfo `json:"current"`
	Previous                PriceInfo `json:"previous"`
	RRP                     PriceInfo `json:"rrp"`
	LowestPriceInLast30Days PriceInfo `json:"lowestPriceInLast30Days"`
	IsMarkedDown            bool      `json:"isMarkedDown"`
	IsOutletPrice           bool      `json:"isOutletPrice"`
	Currency                string    `json:"currency"`
}

type PriceInfo struct {
	Value float64 `json:"value"`
	Text  string  `json:"text"`
}

type FacetGrouping struct {
	Products []FacetProduct `json:"products"`
	Type     string         `json:"type"`
}

type FacetProduct struct {
	ProductID int `json:"productId"`
}
