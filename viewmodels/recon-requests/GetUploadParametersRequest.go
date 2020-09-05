package recon_requests

type GetFileUploadParametersRequest struct {
	SourceFileName string
	SourceFileHash string
	SourceFileColumnCount int
	SourceFileRowCount int

	//comparison file meta data
	ComparisionFileName string
	ComparisonFileHash string
	ComparisonFileColumnCount string
	ComparisonFileRowCount int

	ComparisonPairs []ComparisonPair

}

type ComparisonPair struct {
	SourceColumnIndex int
	ComparisonColumnIndex int
}

func (GetFileUploadParametersRequest) IsValid() error{
	return nil
}
