package validators

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
)

var aValidator *validator.Validate = validator.New()

func ValidateStruct(theStruct recon_requests.ReconRequestsInterface) error {
	err := aValidator.Struct(theStruct)

	//annotation validation failed
	if err != nil {
		return err
	}

	//we can call the custom validation function
	customErr := theStruct.IsValid()

	//custom validation has failed
	if customErr != nil {
		return customErr
	}

	//success
	return nil
}
