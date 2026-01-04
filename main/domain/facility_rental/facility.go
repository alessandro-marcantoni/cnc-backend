package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
)

type Facility struct {
	Id           domain.Id[Facility]
	FacilityType FacilityType
}
