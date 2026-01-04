package servicerental_test

import (
	"testing"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/service_rental"
)

func TestServiceName_String(t *testing.T) {
	tests := []struct {
		name     string
		sn       service_rental.ServiceName
		expected string
	}{
		{"OpenBoardRack", service_rental.OpenBoardRack, "OPEN_BOARD_RACK"},
		{"BoatSpaceInDriftsArea", service_rental.BoatSpaceInDriftsArea, "BOAT_SPACE_IN_DRIFTS_AREA"},
		{"Box", service_rental.Box, "BOX"},
		{"ClosedBoardRack", service_rental.ClosedBoardRack, "CLOSED_BOARD_RACK"},
		{"LargeLocker", service_rental.LargeLocker, "LARGE_LOCKER"},
		{"StandardLocker", service_rental.StandardLocker, "STANDARD_LOCKER"},
		{"SurfBoardStorageWorkshopArea", service_rental.SurfBoardStorageWorkshopArea, "SURF_BOARD_STORAGE_WORKSHOP_AREA"},
		{"ClosedSUPStorage", service_rental.ClosedSUPStorage, "CLOSED_SUP_STORAGE"},
		{"OutdoorCanoeRack", service_rental.OutdoorCanoeRack, "OUTDOOR_CANOE_RACK"},
		{"BoatSpaceVentenana", service_rental.BoatSpaceVentenana, "BOAT_SPACE_VENTENA"},
		{"BoatSpaceTavollo", service_rental.BoatSpaceTavollo, "BOAT_SPACE_TAVOLLO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sn.String(); got != tt.expected {
				t.Errorf("ServiceName.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestToServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected service_rental.ServiceName
	}{
		{"OpenBoardRack", "OPEN_BOARD_RACK", service_rental.OpenBoardRack},
		{"BoatSpaceInDriftsArea", "BOAT_SPACE_IN_DRIFTS_AREA", service_rental.BoatSpaceInDriftsArea},
		{"Box", "BOX", service_rental.Box},
		{"NonExistent", "NON_EXISTENT", service_rental.ServiceName("NON_EXISTENT")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service_rental.ToServiceName(tt.input); got != tt.expected {
				t.Errorf("ToServiceName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServiceName_Equals(t *testing.T) {
	tests := []struct {
		name     string
		sn       service_rental.ServiceName
		other    service_rental.ServiceName
		expected bool
	}{
		{"Same service", service_rental.OpenBoardRack, service_rental.OpenBoardRack, true},
		{"Different services", service_rental.OpenBoardRack, service_rental.Box, false},
		{"Empty services", service_rental.ServiceName(""), service_rental.ServiceName(""), true},
		{"One empty one valid", service_rental.OpenBoardRack, service_rental.ServiceName(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sn.Equals(tt.other); got != tt.expected {
				t.Errorf("ServiceName.Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}
