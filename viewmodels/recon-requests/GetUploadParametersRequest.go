package recon_requests



type GetFileUploadParametersRequest struct {
	SourceFileName string `validate:"required"`
	SourceFileHash string `validate:"required"`
	SourceFileColumnCount int `validate:"gte=1"`
	SourceFileRowCount int `validate:"gte=1"`

	//comparison file meta data
	ComparisionFileName string `validate:"required"`
	ComparisonFileHash string `validate:"required"`
	ComparisonFileColumnCount int `validate:"gte=1"`
	ComparisonFileRowCount int `validate:"gte=1"`

	ComparisonPairs []ComparisonPair `validate:"required,dive,required"`

}

type ComparisonPair struct {
	SourceColumnIndex int `validate:"gte=0"`
	ComparisonColumnIndex int `validate:"gte=0"`
}

func (GetFileUploadParametersRequest) IsValid() error{
	return nil
}
