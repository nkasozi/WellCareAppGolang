package recon_requests

import "gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"

type GetFileUploadParametersRequest struct {
	UserId string `validate:"required"`
	SourceFileName string `validate:"required"`
	SourceFileHash string `validate:"required"`
	SourceFileRowCount int `validate:"gte=1"`

	//comparison file meta data
	ComparisionFileName string `validate:"required"`
	ComparisonFileHash string `validate:"required"`
	ComparisonFileRowCount int `validate:"gte=1"`

	ComparisonPairs []Entities.ComparisonPair `validate:"required,dive,required"`

}


func (GetFileUploadParametersRequest) IsValid() error{
	return nil
}
