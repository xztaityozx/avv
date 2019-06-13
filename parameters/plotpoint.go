package parameters

type PlotPoint struct {
	Start   float64
	Step    float64
	Stop    float64
	Signals []string
}

// Signal Name and filter values for count up failure stage
type Filter struct {
	SignalName string
	Status     []string
}
