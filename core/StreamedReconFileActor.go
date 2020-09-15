package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core/gcloud"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"time"
)

func ProcessStreamReconFileApiRequest(req recon_requests.StreamReconRequest) (processingResult *recon_responses.StreamReconFileResponse, err error) {

	defer func() (*recon_responses.StreamReconFileResponse, error) {

		//catch any exceptions and return a clean error
		if r := recover(); r != nil {

			err = errors.New(fmt.Sprint("Error:", r))
			return nil, err
		}

		//success response return as is
		return processingResult, err
	}()

	//get details of the upload request
	//those details should have the topics
	originalUploadRequest, err := redisDataStore.GetById(req.UploadRequestId)

	if err != nil {
		return nil, err
	}

	switch req.FileType {

	case UploadFileTypes.SrcFile:
		return processStreamSourceFileReconRequest(req, originalUploadRequest)

	case UploadFileTypes.ComparisonFile:
		return processStreamComparisonFileReconRequest(req, originalUploadRequest)

	default:
		return processStreamUnknownFileReconRequest()
	}
}

//src file topic will be pull based subscription
//comparison file topic will be push based
func processStreamSourceFileReconRequest(req recon_requests.StreamReconRequest, originalUploadDetails Entities.FilesUploadedParameters) (*recon_responses.StreamReconFileResponse, error) {

	streamedFileChunk := Entities.StreamedFileChunk{
		UploadRequestId:     req.UploadRequestId,
		FileType:            req.FileType,
		ChunkSequenceNumber: req.ChunkSequenceNumber,
		Records:             req.Records,
		Id:                  shared.GenerateUniqueId("SRC-FILE-CHUNK-"),
		DateCreated:         time.Now(),
		DateModified:        time.Now(),
	}

	jsonBytes, err := json.Marshal(&streamedFileChunk)

	if err != nil {
		return nil, err
	}

	topicForThisSourceFile := GenerateNameForSourceFileTopic(originalUploadDetails.UserId, originalUploadDetails.SourceFileHash, req.ChunkSequenceNumber)

	messageId, err := gcloud.PublishMessageToCloudRunTopic(jsonBytes, Constants.GOOGLE_CLOUD_PROJECT_ID, topicForThisSourceFile)

	if err != nil {
		return nil, err
	}

	processingResult := recon_responses.StreamReconFileResponse{
		MessageId: messageId,
		Status:    "OK",
	}

	return &processingResult, nil
}



func processStreamComparisonFileReconRequest(req recon_requests.StreamReconRequest, originalUploadDetails Entities.FilesUploadedParameters) (*recon_responses.StreamReconFileResponse, error) {
	streamedFileChunk := Entities.StreamedFileChunk{
		UploadRequestId:     req.UploadRequestId,
		FileType:            req.FileType,
		ChunkSequenceNumber: req.ChunkSequenceNumber,
		Records:             req.Records,
		Id:                  shared.GenerateUniqueId("CMP-FILE-CHUNK-"),
		DateCreated:         time.Now(),
		DateModified:        time.Now(),
	}

	jsonBytes, err := json.Marshal(&streamedFileChunk)

	if err != nil {
		return nil, err
	}

	topicForThisComparisonFile := GenerateNameForComparisonFileTopic(originalUploadDetails.UserId, originalUploadDetails.ComparisonFileHash)

	messageId, err := gcloud.PublishMessageToCloudRunTopic(jsonBytes, Constants.GOOGLE_CLOUD_PROJECT_ID, topicForThisComparisonFile)

	if err != nil {
		return nil, err
	}

	processingResult := recon_responses.StreamReconFileResponse{
		MessageId: messageId,
		Status:    "OK",
	}

	return &processingResult, nil
}

func processStreamUnknownFileReconRequest() (*recon_responses.StreamReconFileResponse, error) {
	processingResult := recon_responses.StreamReconFileResponse{}

	return &processingResult, errors.New("Unknown File Type Supplied")
}
