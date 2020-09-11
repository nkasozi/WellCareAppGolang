package core

import (
	"errors"
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"time"
)

var redisDataStore = datastore.RedisDataStoreFilesUplaodedParameters{}

//determines the batch for the file upload
func determineBatchSizeForFile(req recon_requests.GetFileUploadParametersRequest, fileType UploadFileTypes.UploadFileType) int {
	switch fileType {
	case UploadFileTypes.ComparisonFile:
		return 200
	case UploadFileTypes.SrcFile:
		return 200
	default:
		return 200
	}
}

//saves an entity in a datastore
func saveUploadRequest(req Entities.FilesUploadedParameters) (string, error) {

	//create request
	if req.Id == "0" {
		return redisDataStore.Add(req)
	}

	//update request
	return redisDataStore.Update(req)
}

func ProcessGetUploadParametersRequest(req recon_requests.GetFileUploadParametersRequest) (processingResult recon_responses.GetFileUploadParametersResponse, err error) {

	processingResult = recon_responses.GetFileUploadParametersResponse{}

	defer func() (recon_responses.GetFileUploadParametersResponse, error) {

		//catch any exceptions and return a clean error
		if r := recover(); r != nil {

			err = errors.New(fmt.Sprint("Error:", r))
			return processingResult, err
		}

		//success response return as is
		return processingResult, err
	}()

	//check if we have ever recieved these files before
	response, checkErr := checkIfThisIsARepeatUploadAttempt(req, processingResult)

	//yes its a repeat attempt
	if checkErr == nil {
		return response, nil
	}

	//process the request
	return processNewFileUpload(req, processingResult)
}

func checkIfThisIsARepeatUploadAttempt(req recon_requests.GetFileUploadParametersRequest, processingResult recon_responses.GetFileUploadParametersResponse) (recon_responses.GetFileUploadParametersResponse, error) {

	//check if we have recived this item before
	fileId := generateFileUploadID(req)
	existingFileUploadAttempt, getErr := redisDataStore.GetById(fileId)

	//error in checking for item in redis
	if getErr != nil {
		return processingResult, getErr
	}

	//means we have no existing attempt
	if existingFileUploadAttempt.Id == "" {
		return processingResult, errors.New("No previous Attempt has been logged for File with ID ["+fileId+"]")
	}

	//build up response
	processingResult.SourceFileExpectedBatchSize = existingFileUploadAttempt.SourceFileExpectedBatchSize
	processingResult.SourceFileHash = req.SourceFileHash
	processingResult.SourceFileName = req.SourceFileName
	processingResult.ComparisonFileHash = req.ComparisonFileHash
	processingResult.ComparisonFileName = req.ComparisionFileName
	processingResult.ComparisonFileExpectedBatchSize = existingFileUploadAttempt.ComparisonFileExpectedBatchSize
	processingResult.UploadRequestId = existingFileUploadAttempt.Id
	processingResult.SourceFileIsFirstTimeUpload = false
	processingResult.SourceFileLastRowReceived = existingFileUploadAttempt.SourceFileLastRowReceived
	processingResult.ComparisonFileIsFirstTimeUpload = false
	processingResult.ComparisonFileLastRowReceived = existingFileUploadAttempt.ComparisonFileLastRowReceived

	//sucess
	return processingResult, nil
}

func processNewFileUpload(req recon_requests.GetFileUploadParametersRequest, processingResult recon_responses.GetFileUploadParametersResponse) (recon_responses.GetFileUploadParametersResponse, error) {

	//determine the batch size for each upload request
	srcFileBatchSize := determineBatchSizeForFile(req, UploadFileTypes.SrcFile)
	comparisonFileBatchSize := determineBatchSizeForFile(req, UploadFileTypes.ComparisonFile)

	//map to datatstore struct to save
	filesUploadedParameters := Entities.FilesUploadedParameters{
		UserId:                          req.UserId,
		SourceFileName:                  req.SourceFileName,
		SourceFileHash:                  req.SourceFileHash,
		SourceFileRowCount:              req.SourceFileRowCount,
		SourceFileExpectedBatchSize:     srcFileBatchSize,
		SourceFileLastRowReceived:       0,
		ComparisionFileName:             req.ComparisionFileName,
		ComparisonFileHash:              req.ComparisonFileHash,
		ComparisonFileRowCount:          req.ComparisonFileRowCount,
		ComparisonFileExpectedBatchSize: comparisonFileBatchSize,
		ComparisonFileLastRowReceived:   0,
		ComparisonPairs:                 req.ComparisonPairs,
		Id:                              generateFileUploadID(req),
		DateCreated:                     time.Now(),
		DateModified:                    time.Now(),
	}

	//save the upload request meta data in redis
	uploadRequestId, saveErr := saveUploadRequest(filesUploadedParameters)

	//failed to save..stop here
	if saveErr != nil {
		return processingResult, saveErr
	}

	//build up response
	processingResult.SourceFileExpectedBatchSize = srcFileBatchSize
	processingResult.SourceFileHash = req.SourceFileHash
	processingResult.SourceFileName = req.SourceFileName
	processingResult.ComparisonFileHash = req.ComparisonFileHash
	processingResult.ComparisonFileName = req.ComparisionFileName
	processingResult.ComparisonFileExpectedBatchSize = comparisonFileBatchSize
	processingResult.UploadRequestId = uploadRequestId
	processingResult.SourceFileIsFirstTimeUpload = true
	processingResult.SourceFileLastRowReceived = 0
	processingResult.ComparisonFileIsFirstTimeUpload = true
	processingResult.ComparisonFileLastRowReceived = 0

	//done..success
	return processingResult, nil
}

func generateFileUploadID(req recon_requests.GetFileUploadParametersRequest) string {
	return req.UserId + "-" + req.SourceFileHash + "-" + req.ComparisonFileHash
}
