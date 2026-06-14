package costing

// CostOptions configures unit rates for cost calculation.
type CostOptions struct {
	Currency           string
	MaterialRatesM2    map[string]float64
	EdgeBandingRatesM  map[string]float64
	HardwareRates      map[string]float64
	LaborRatePerHour   float64
	WastePercentage    float64
}

func DefaultCostOptions() CostOptions {
	return CostOptions{
		Currency: "USD",
		MaterialRatesM2: map[string]float64{
			"melamine-white-18": 38,
			"melamine-default-18": 35,
			"melamine-white-18-back": 12,
		},
		EdgeBandingRatesM: map[string]float64{
			"pvc-white-1mm": 1.2,
		},
		HardwareRates: map[string]float64{
			"hinge":          4.5,
			"drawer_runner":  28,
			"hanger_rod":     18,
			"rod_bracket":    2.5,
		},
		LaborRatePerHour: 28,
		WastePercentage:  0.12,
	}
}

func (o CostOptions) materialRate(materialID string) float64 {
	if rate, ok := o.MaterialRatesM2[materialID]; ok {
		return rate
	}
	return o.MaterialRatesM2["melamine-default-18"]
}

func (o CostOptions) edgeBandingRate(material string) float64 {
	if rate, ok := o.EdgeBandingRatesM[material]; ok {
		return rate
	}
	return 1.0
}

func (o CostOptions) hardwareRate(hwType string) float64 {
	if rate, ok := o.HardwareRates[hwType]; ok {
		return rate
	}
	return 5.0
}
