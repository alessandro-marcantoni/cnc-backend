package facilityrental

type FacilityName string

const (
	OpenBoardRack                FacilityName = "OPEN_BOARD_RACK"
	BoatSpaceInDriftsArea        FacilityName = "BOAT_SPACE_IN_DRIFTS_AREA"
	Box                          FacilityName = "BOX"
	ClosedBoardRack              FacilityName = "CLOSED_BOARD_RACK"
	LargeLocker                  FacilityName = "LARGE_LOCKER"
	StandardLocker               FacilityName = "STANDARD_LOCKER"
	SurfBoardStorageWorkshopArea FacilityName = "SURF_BOARD_STORAGE_WORKSHOP_AREA"
	ClosedSUPStorage             FacilityName = "CLOSED_SUP_STORAGE"
	OutdoorCanoeRack             FacilityName = "OUTDOOR_CANOE_RACK"
	BoatSpaceVentenana           FacilityName = "BOAT_SPACE_VENTENA"
	BoatSpaceTavollo             FacilityName = "BOAT_SPACE_TAVOLLO"
)

func (sn FacilityName) String() string {
	return string(sn)
}

func ToFacilityName(s string) FacilityName {
	return FacilityName(s)
}

func (sn FacilityName) Equals(other FacilityName) bool {
	return sn == other
}
