package structs

type Estate struct {
	ID                  int
	PropertyDescription PropertyDescription
	Interior            InteriorFeatures
	Exterior            ExteriorFeatures
	OtherDetails        OtherDetails
	Status              Status
	Luxury              bool
	Images              []string

	CreateAt string
	UpdateAt string
}

type EstatesResponse struct {
	Estates []EstateForList
	Total   int
}

type EstateForList struct {
	ID          int
	Name        string
	Price       int
	Address     string
	Beds        int
	Baths       int
	AreaInMeter int
	Latitude    string
	Longitude   string
	Images      []string
	Status      Status

	//Country     int
	//City        int
}

type PropertyDescription struct {
	Name        string
	Price       int
	Country     int
	City        int
	Address     string
	Beds        int
	Baths       int
	AreaInMeter int
	Type        PropertyType
	YearBuilt   int
	Description string
	Latitude    string
	Longitude   string
}

type InteriorFeatures struct {
	Appliances       string
	Features         string
	KitchenFeatures  string
	TotalBedrooms    int
	FullBathrooms    int
	HalfBathrooms    int
	FloorDescription string
	Fireplace        string
	Cooling          int // -- bool 0 - null, 1 - false, 2 - true
	Heating          int // -- bool 0 - null, 1 - false, 2 - true
}

type ExteriorFeatures struct {
	LotSizeInAcres int
	Features       string
	ArchStyle      string
	Roof           string
	Sewer          string
}

type OtherDetails struct {
	AreaName        string
	Garage          int
	Parking         string
	View            string
	Pool            int // -- bool 0 - null, 1 - false, 2 - true
	PoolDescription string
	WaterSource     string
	Utilities       string
}

type PropertyType int

type PropertyTypeKeyValue struct {
	Key   int
	Value string
}

type Status string
