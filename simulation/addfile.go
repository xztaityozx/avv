package simulation

// AddFile Addfileを表すstruct
type AddFile struct {
	// VddVoltage Vddの電圧
	VddVoltage float64
	// GndVoltage Gndの電圧
	GndVoltage float64
	// ICCommand
	ICCommand string
	// オプション
	Options []string
	// SEED Seed値
	SEED int
}
