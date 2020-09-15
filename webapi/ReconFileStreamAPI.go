package webapi

import (
	"encoding/json"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/validators"
	"net/http"
)

func StreamReconFile(response http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodOptions:
		//handle CORs preflight request
		shared.GenerateCORsResponse(response, http.MethodPost)
		return

	case http.MethodPost:
		//handle POST /StreamReconFile
		processStreamReconFilePostRequest(response, request)
		return

	default:
		//oops unhandled method invoked
		shared.GenerateMethodNotAllowedResponse(response, request.Method)
		return
	}
}

func processStreamReconFilePostRequest(response http.ResponseWriter, request *http.Request) {

	//decode request to view model
	var decodedRequest = recon_requests.StreamReconRequest{}
	err := json.NewDecoder(request.Body).Decode(&decodedRequest)

	//failed to decode request
	if err != nil {
		shared.GenerateBadRequestResponse(response, err)
		return
	}

	//validate the view model
	if err := validators.ValidateStruct(decodedRequest); err != nil {
		shared.GenerateBadRequestResponse(response, err)
		return
	}

	//process the request
	processingResult, processingErr := core.ProcessStreamReconFileApiRequest(decodedRequest)

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
