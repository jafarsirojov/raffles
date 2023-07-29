package structs

type SearchOptions struct {
	PriceMin         int
	PriceMax         int
	BedsMax          int
	BathsMax         int
	PropertyTypes    []PropertyTypeKeyValue
	SquareFootageMin int
	SquareFootageMax int
	LotSizeMin       int
	LotSizeMax       int
	YearBuiltMin     int
	YearBuiltMax     int
	GarageSpacesMin  int
	GarageSpacesMax  int
}
