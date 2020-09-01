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
}

func (GetFileUploadParametersRequest) IsValid() error{
	return nil
}
