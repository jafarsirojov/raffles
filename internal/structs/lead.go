package structs

type Lead struct {
	ID                     int
	Name                   string
	Site                   string
	REStageConstruction    string
	RERegion               string
	REType                 string
	REPurposeOfAcquisition string
	RECountOfRooms         string
	PurchaseBudget         string
	Phone                  string
	Email                  string
	CommunicationMethod    string
	Description            string
	Status                 string

	CreatedAt string
	UpdatedAt string
}

type Comment struct {
	ID         int
	LeadID     int
	AdminID    int
	AdminLogin string
	Text       string
	CreatedAt  string
}

type LeadAndComments struct {
	Lead     Lead
	Comments []Comment
}
