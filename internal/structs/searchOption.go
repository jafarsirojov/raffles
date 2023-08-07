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

	Status []Status
}

type SearchOptionsDTO struct {
	PriceMin         int
	PriceMax         int
	BedsMin          int
	BathsMin         int
	PropertyTypes    []int
	SquareFootageMin int
	SquareFootageMax int
	LotSizeMin       int
	LotSizeMax       int
	YearBuiltMin     int
	YearBuiltMax     int
	GarageSpacesMin  int
	GarageSpacesMax  int
	Cooling          int // -- bool 0 - null, 1 - false, 2 - true
	Heating          int // -- bool 0 - null, 1 - false, 2 - true
	Pool             int // -- bool 0 - null, 1 - false, 2 - true
}
