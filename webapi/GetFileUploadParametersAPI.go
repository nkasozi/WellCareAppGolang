package webapi

import (
	"encoding/json"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon-requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/validators"
	"net/http"
)



func processGetUploadParametersRequest(req recon_requests.GetFileUploadParametersRequest) (recon_responses.GetFileUploadParametersResponse, error) {

	response := recon_responses.GetFileUploadParametersResponse{}

	//process the request
	//TODO

	//build up response
	response.SourceFileExpectedBatchSize = 200
	response.SourceFileHash = req.SourceFileHash
	response.SourceFileName = req.SourceFileName
	response.ComparisonFileHash = req.ComparisonFileHash
	response.ComparisonFileName = req.ComparisionFileName
	response.ComparisonFileExpectedBatchSize = 200

	return response, nil
}

func processHttpPostRequest(response http.ResponseWriter, request *http.Request) {

	//decode request to view model
	var decodedRequest = recon_requests.GetFileUploadParametersRequest{}
	err := json.NewDecoder(request.Body).Decode(&decodedRequest)

	//failed to decode request
	if err != nil {
		shared.GenerateBadRequestResponse(response, err)
		return
	}

	//validate the view model
	if err := validators.ValidateStruct(decodedRequest); err != nil {
		shared.GenerateBadRequestResponse(response,err)
		return
	}

	//process the request
	processingResult, processingErr := processGetUploadParametersRequest(decodedRequest)

	//something serious happened when processing the request
	if processingErr != nil {
		shared.GenerateInternalServerResponse(response, processingErr)
		return
	}

	//get the json version of the struct
	jsonString, encodingErr := shared.ToJsonString(processingResult)

	//failed to encode it as json
	if encodingErr != nil {
		shared.GenerateInternalServerResponse(response, encodingErr)
		return
	}

	//return result of processing
	shared.GenerateOkRequestResponse(response, jsonString)
}

// PingExample godoc
// @Summary ping example
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /examples/ping [get]
func GetFileUploadParameters(response http.ResponseWriter, request *http.Request) {

	switch request.Method {


	case http.MethodOptions:
		//handle CORs preflight request
		shared.GenerateCORsResponse(response, http.MethodPost)
		return

	case http.MethodPost:
		//handle POST /GetFileUploadParameters
		processHttpPostRequest(response, request)
		return

	default:
		//oops unhandled method invoked
		shared.GenerateMethodNotAllowedResponse(response, request.Method)
		return
	}

}
