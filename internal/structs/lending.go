package structs

type Lending struct {
	ID                      int
	Name                    string
	FullName                string
	Address                 string
	StartingPrice           MultiCurrency
	ListingDetails          ListingDetails
	FeaturesAndAmenities    []FeatureOrAmenity
	FeaturesAndAmenitiesIDs []int
	Title                   string
	Description             string
	Video                   string
	Images                  []string
	Availabilities          []Availability
	CreatedAt               string
	UpdatedAt               string
}

type LendingData struct {
	Name                 string
	FullName             string
	Address              string
	StartingPrice        MultiCurrency
	ListingDetails       ListingDetails
	FeaturesAndAmenities []FeatureOrAmenity
	Title                string
	Description          string
	Video                string
	Images               []string
	Availabilities       []Availability
	ImageURL             string
}

type LendingList struct {
	ID   int
	Name string
}

type Availability struct {
	ID          int
	LendingID   int
	Price       MultiCurrency
	UniqueID    string
	Bedroom     int
	Parking     int
	Area        float32
	Plot        float32
	SpecialGift string
}

type MultiCurrency struct {
	AED int
	USD int
}

type ListingDetails struct {
	PropertyType string
	Furnishing   string
}

type FeatureOrAmenity struct {
	ID   int
	Icon string
	Name string
}
