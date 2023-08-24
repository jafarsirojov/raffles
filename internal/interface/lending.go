package interfaces

import (
	"context"
	"crm/internal/structs"
)

type LendingAdminRepo interface {
	SaveLending(ctx context.Context, data structs.Lending) error
	UpdateLending(ctx context.Context, data structs.Lending) (err error)
	GetLendingList(ctx context.Context) (list []structs.LendingList, err error)
	GetLendingByID(ctx context.Context, id int) (data structs.Lending, err error)

	GetImagesByLendingID(ctx context.Context, id int) (images []string, err error)
	UpdateLendingImages(ctx context.Context, id int, images []string) error
	SelectBackgroundImageLandingID(ctx context.Context, id int) (backgroundImage string, err error)
	UpdateBackgroundImage(ctx context.Context, id int, new string) error
	SelectMainLogoByLandingID(ctx context.Context, id int) (logo string, err error)
	UpdateMainLogo(ctx context.Context, id int, new string) error
	SelectPartnerLogoByLandingID(ctx context.Context, id int) (logo string, err error)
	UpdatePartnerLogo(ctx context.Context, id int, new string) error
	SelectOurLogoByLandingID(ctx context.Context, id int) (logo string, err error)
	UpdateOurLogo(ctx context.Context, id int, new string) error

	SelectFeaturesAndAmenities(ctx context.Context) (list []structs.FeatureOrAmenity, err error)
	SelectFeaturesAndAmenitiesByIDs(ctx context.Context, ids []int) (list []structs.FeatureOrAmenity, err error)
	InsertFeatureAndAmenity(ctx context.Context, name, icon string) error
	DeleteFeatureAndAmenity(ctx context.Context, id int) error

	SaveAvailability(ctx context.Context, data structs.Availability) (err error)
	GetAvailabilitiesByLandingID(ctx context.Context, landingID int) (list []structs.Availability, err error)
	UpdateAvailability(ctx context.Context, data structs.Availability) (err error)
	DeleteAvailability(ctx context.Context, id int) (err error)
	SelectPaymentPlanByAvailabilityID(ctx context.Context, id int) (paymentPlan string, err error)
	UpdatePaymentPlan(ctx context.Context, id int, new string) error
}

type LendingClientRepo interface {
	SelectLandingList(ctx context.Context) (list []structs.LendingListMainPage, err error)
	SelectLendingData(ctx context.Context, id int) (data structs.Lending, err error)
	SelectFeaturesAndAmenitiesByIDs(ctx context.Context, ids []int) (list []structs.FeatureOrAmenity, err error)
	GetAvailabilitiesByLandingID(ctx context.Context, landingID int) (list []structs.Availability, err error)
}
