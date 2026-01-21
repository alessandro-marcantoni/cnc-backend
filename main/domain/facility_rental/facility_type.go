package facilityrental

import "github.com/alessandro-marcantoni/cnc-backend/main/domain"

type FacilityType struct {
	Id             domain.Id[FacilityType]
	FacilityName   FacilityName
	Description    string
	SuggestedPrice float64
	HasBoat        bool
}
