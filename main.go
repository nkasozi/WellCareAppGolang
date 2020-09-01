package backend_core

import (
	"encoding/json"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon-requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"log"
	"net/http"
)

func processGetUploadParametersRequest(req recon_requests.GetFileUploadParametersRequest) (recon_responses.GetFileUploadParametersResponse, error) {

	response := recon_responses.GetFileUploadParametersResponse{}

	//validate the view model
	if err := req.IsValid(); err != nil {
		return response, err
	}

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

func GetFileUploadParameters(response http.ResponseWriter, request *http.Request) {

	//decode request to view model
	var decodedRequest = recon_requests.GetFileUploadParametersRequest{}
	err := json.NewDecoder(request.Body).Decode(&decodedRequest)

	//failed to decode request
	if err != nil {
		generateBadRequestResponse(response, err)
		return
	}

	//process the request
	processingResult, processingErr := processGetUploadParametersRequest(decodedRequest)

	if processingErr != nil {
		generateInternalServerResponse(response, processingErr)
		return
	}

	//get the json version of the struct
	jsonString, encodingErr := shared.ToJsonString(processingResult)

	//failed to encode it as json
	if encodingErr != nil {
		generateInternalServerResponse(response, encodingErr)
		return
	}

	//return result of processing
	generateOkRequestResponse(response, jsonString)
}

func generateBadRequestResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type","text/plain")
	w.Write([]byte(err.Error()))
	log.Println(err.Error())
	return
}

func generateOkRequestResponse(w http.ResponseWriter, json string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	w.Write([]byte(json))
}

func generateInternalServerResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type","text/plain")
	w.Write([]byte(err.Error()))
	log.Println(err.Error())
	return
}
