package facilityrental

type BoatInfo struct {
	Name          string
	LengthMeters  float64
	WidthMeters   *float64 // Nullable - can be nil if not measured
	Type          string   // Type/category of boat (e.g., Sailing, Motor, Inflatable)
	EngineInfo    string
	InsuranceInfo BoatInsuranceInfo
}

func (b BoatInfo) HasInsurance() bool {
	return b.InsuranceInfo.HasInsurance()
}

type BoatInsuranceInfo interface {
	HasInsurance() bool
}

type BoatInsurance struct {
	ProviderName   string
	PolicyNumber   string
	ExpirationDate string
}

type NoBoatInsurance struct{}

func (b BoatInsurance) HasInsurance() bool {
	return true
}

func (n NoBoatInsurance) HasInsurance() bool {
	return false
}
