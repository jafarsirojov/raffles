package interfaces

import (
	"context"
	"crm/internal/structs"
)

type LandingAdminRepo interface {
	SaveLanding(ctx context.Context, data structs.Landing) error
	UpdateLanding(ctx context.Context, data structs.Landing) (err error)
	GetLandingList(ctx context.Context, offset, limit int) (list []structs.LandingList, count int, err error)
	GetLandingByID(ctx context.Context, id int) (data structs.Landing, err error)

	GetImagesByLandingID(ctx context.Context, id int) (images []string, err error)
	UpdateLandingImages(ctx context.Context, id int, images []string) error
	SelectBackgroundImageByLandingID(ctx context.Context, id int) (backgroundImage string, err error)
	UpdateBackgroundImage(ctx context.Context, id int, new string) error
	SelectBackgroundForMobileByLandingID(ctx context.Context, id int) (backgroundImage string, err error)
	UpdateBackgroundForMobile(ctx context.Context, id int, new string) error
	SelectMainLogoByLandingID(ctx context.Context, id int) (logo string, err error)
	UpdateMainLogo(ctx context.Context, id int, new string) error
	SelectPartnerLogoByLandingID(ctx context.Context, id int) (logo string, err error)
	UpdatePartnerLogo(ctx context.Context, id int, new string) error
	SelectOurLogoByLandingID(ctx context.Context, id int) (logo string, err error)
	UpdateOurLogo(ctx context.Context, id int, new string) error
	SelectVideoCoverByLandingID(ctx context.Context, id int) (cover string, err error)
	UpdateVideoCover(ctx context.Context, id int, new string) error

	SelectFeaturesAndAmenities(ctx context.Context) (list []structs.FeatureOrAmenity, err error)
	SelectFeaturesAndAmenitiesByIDs(ctx context.Context, ids []int) (list []structs.FeatureOrAmenity, err error)
	InsertFeatureAndAmenity(ctx context.Context, name, icon string) error
	DeleteFeatureAndAmenity(ctx context.Context, id int) error

	SaveAvailability(ctx context.Context, data structs.Availability) (err error)
	GetAvailabilitiesByLandingID(ctx context.Context, landingID int) (list []structs.Availability, err error)
	UpdateAvailability(ctx context.Context, data structs.Availability) (err error)
	DeleteAvailability(ctx context.Context, id int) (err error)
	SelectFilePlanByLandingID(ctx context.Context, id int) (paymentPlan string, err error)
	UpdateFilePlan(ctx context.Context, id int, new string) error
}

type LandingClientRepo interface {
	SelectLandingList(ctx context.Context) (list []structs.LandingListMainPage, err error)
	SelectLandingData(ctx context.Context, id int) (data structs.Landing, err error)
	SelectFeaturesAndAmenitiesByIDs(ctx context.Context, ids []int) (list []structs.FeatureOrAmenity, err error)
	GetAvailabilitiesByLandingID(ctx context.Context, landingID int) (list []structs.Availability, err error)
}
