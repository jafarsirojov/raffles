package structs

type Lending struct {
	ID                      int
	Name                    string
	MainDescription         string
	FullName                string
	Slogan                  string
	Address                 string
	StartingPrice           MultiCurrency
	ListingDetails          ListingDetails
	FeaturesAndAmenities    []FeatureOrAmenity
	FeaturesAndAmenitiesIDs []int
	Title                   string
	Description             string
	Video                   string
	FilePlan                string
	TitlePlan               string
	BackgroundImage         string
	MainLogo                string
	PartnerLogo             string
	OurLogo                 string
	Images                  []string
	Availabilities          []Availability
	CreatedAt               string
	UpdatedAt               string
}

type LendingData struct {
	Name                 string
	MainDescription      string
	FullName             string
	Slogan               string
	Address              string
	StartingPrice        MultiCurrency
	ListingDetails       ListingDetails
	FeaturesAndAmenities []FeatureOrAmenity
	Title                string
	Description          string
	Video                string
	FilePlan             string
	TitlePlan            string
	BackgroundImage      string
	MainLogo             string
	PartnerLogo          string
	OurLogo              string
	Images               []string
	Availabilities       []Availability
	FileURL              string
}

type LendingList struct {
	ID   int
	Name string
}

type ListMainPage struct {
	List    []LendingListMainPage
	FileURL string
}

type LendingListMainPage struct {
	Name            string
	MainDescription string
	BackgroundImage string
}

type Availability struct {
	ID          int
	LendingID   int
	Price       MultiCurrency
	UniqueID    string
	Bedroom     int
	Parking     int
	Area        string
	Plot        string
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
