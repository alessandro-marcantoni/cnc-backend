package facilityrental

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type RentalManagementService struct {
	repository FacilityRepository
}

func NewRentalManagementService(repository FacilityRepository) *RentalManagementService {
	return &RentalManagementService{repository: repository}
}

func (this RentalManagementService) RentService(f Facility, year time.Time) result.Result[RentedFacility] {
	availableFacilities := this.repository.GetAvailableFacilities(f.FacilityType.FacilityName)
	if len(availableFacilities) == 0 {
		return result.Err[RentedFacility](errors.RentError{Description: "no available services of this type"})
	}
	return result.Err[RentedFacility](errors.RentError{Description: "not implemented"})
}

func (this RentalManagementService) GetFacilitiesCatalog() []FacilityType {
	return this.repository.GetFacilitiesCatalog()
}

func (this RentalManagementService) GetFacilitiesByType(facilityTypeId domain.Id[FacilityType]) []FacilityWithStatus {
	return this.repository.GetFacilitiesByType(facilityTypeId)
}
