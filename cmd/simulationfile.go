package cmd

type Templates struct {
	AddFile       string
	SPIFile       string
	ResultsXML    string
	ResultsMapXML string
}

// Make AddFile
// return [add file path]
//func (s SEED) MkAddFile(base string) []string {
//	format := FU.Cat(config.Templates.AddFile)
//	var rt []string
//
//	for seed := s.Start; seed <= s.End; seed++ {
//		path := PathJoin(base, "AddFiles", fmt.Sprintf("SEED%05d", seed))
//		rt = append(rt, path)
//
//	}
//
//}

type SimulationFiles struct {
	AddFile       string
	SPIScript     string
	ACEScript     string
	ResultsXML    string
	ResultsMapXML string
	ModelFile     string
}

type SimulationDirectories struct {
	DstDir     string
	NetListDir string
	BaseDir    string
}

// compare func for SimulationFiles struct
//returns: compare result
func (s SimulationFiles) Compare(t SimulationFiles) bool {
	return s.ACEScript == t.ACEScript &&
		s.ResultsMapXML == t.ResultsMapXML &&
		t.ResultsXML == s.ResultsXML &&
		s.SPIScript == t.SPIScript &&
		s.AddFile == t.AddFile &&
		s.ModelFile == t.ModelFile
}
