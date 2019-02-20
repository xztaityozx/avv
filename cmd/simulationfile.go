package cmd

type Templates struct {
	SPIScript     string
	ResultsMapXML string
}

type SimulationFiles struct {
	AddFile       AddFile
	SPIScript     string
	ACEScript     string
	ResultsXML    string
	ResultsMapXML string
	ModelFile     string
	Self          string
}

type SimulationDirectories struct {
	DstDir     string
	NetListDir string
	BaseDir    string
	SearchDir  string
	ResultDir  string
}

// compare func for SimulationFiles struct
//returns: compare result
func (s SimulationFiles) Compare(t SimulationFiles) bool {
	return s.ACEScript == t.ACEScript &&
		s.ResultsMapXML == t.ResultsMapXML &&
		t.ResultsXML == s.ResultsXML &&
		s.SPIScript == t.SPIScript &&
		s.AddFile.Compare(t.AddFile) &&
		s.ModelFile == t.ModelFile
}
