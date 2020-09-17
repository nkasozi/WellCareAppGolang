package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core/pubsubserviceproviders"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core/pubsubserviceproviders/gcloud"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"time"
)

var srcFileTopicNotFoundErr = errors.New("Unable to find Src file topic for this chunk sequence Number")
var PubSubServiceProvider pubsubserviceproviders.PubSubInterface = gcloud.NewPubSubClient()

func ProcessStreamReconFileApiRequest(req recon_requests.StreamFileChunkForReconRequest) (processingResult *recon_responses.StreamFileChunkForReconResponse, err error) {

	//make sure we catch any panics etc
	defer func() (*recon_responses.StreamFileChunkForReconResponse, error) {

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

	//cant find original
	//upload request
	if err != nil {
		return nil, err
	}

	//handle differently depending
	//on file type
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
func processStreamSourceFileReconRequest(req recon_requests.StreamFileChunkForReconRequest, originalUploadDetails Entities.FilesUploadedParameters) (*recon_responses.StreamFileChunkForReconResponse, error) {

	//search for the right topic for this chunk
	topicForThisSourceFile, err := getTopicForThisSrcFileChunk(originalUploadDetails.SourceFileTopics, req.ChunkSequenceNumber)

	//oops cant find the topic meant
	//for this chunk
	if err != nil {
		return nil, err
	}

	//we can now build the job THAT will be published
	//to the QueueService
	streamedFileChunk := Entities.StreamedFileChunk{
		UploadRequestId:     req.UploadRequestId,
		FileType:            req.FileType,
		ChunkSequenceNumber: req.ChunkSequenceNumber,
		Records:             req.Records,
		IsEOF:               req.IsEOF,
		Id:                  shared.GenerateUniqueId("SRC-FILE-CHUNK-"),
		DateCreated:         time.Now(),
		DateModified:        time.Now(),
	}

	//turn it to json
	jsonBytes, err := json.Marshal(&streamedFileChunk)

	//oops error on trying to change
	//it to Json
	if err != nil {
		return nil, err
	}

	//publish the json to the queue servic
	messageId, err := PubSubServiceProvider.PublishMessageToPubSubTopic(jsonBytes, Constants.GOOGLE_CLOUD_PROJECT_ID, topicForThisSourceFile.TopicName)

	//oops error
	//on publish
	if err != nil {
		return nil, err
	}

	//by this time we have sucess
	processingResult := recon_responses.StreamFileChunkForReconResponse{
		MessageId:                   messageId,
		Status:                      "OK",
		OriginalChunkSequenceNumber: req.ChunkSequenceNumber,
		OriginalFileType:            req.FileType,
	}

	//return result
	return &processingResult, nil
}

func getTopicForThisSrcFileChunk(topics []Entities.Topic, number int) (*Entities.Topic, error) {

	for _, topic := range topics {

		if topic.ChunkSequenceNumber == number {
			return &topic, nil
		}

	}

	return nil, srcFileTopicNotFoundErr
}

func processStreamComparisonFileReconRequest(req recon_requests.StreamFileChunkForReconRequest, originalUploadDetails Entities.FilesUploadedParameters) (*recon_responses.StreamFileChunkForReconResponse, error) {

	//we can now build the message that  will be published
	//to the QueueService
	streamedFileChunk := Entities.StreamedFileChunk{
		UploadRequestId:     req.UploadRequestId,
		FileType:            req.FileType,
		ChunkSequenceNumber: req.ChunkSequenceNumber,
		Records:             req.Records,
		IsEOF:               req.IsEOF,
		Id:                  shared.GenerateUniqueId("CMP-FILE-CHUNK-"),
		DateCreated:         time.Now(),
		DateModified:        time.Now(),
	}

	//we turn the request to json
	jsonBytes, err := json.Marshal(&streamedFileChunk)

	//oops error on jsonifying
	if err != nil {
		return nil, err
	}

	//get topic for the comparison file
	topicForThisComparisonFile := originalUploadDetails.ComparisonFileTopic

	//publish the message to the topic
	messageId, err := PubSubServiceProvider.PublishMessageToPubSubTopic(jsonBytes, Constants.GOOGLE_CLOUD_PROJECT_ID, topicForThisComparisonFile.TopicName)

	//err on publish
	if err != nil {
		return nil, err
	}

	//uild success response
	processingResult := recon_responses.StreamFileChunkForReconResponse{
		MessageId:                   messageId,
		Status:                      "OK",
		OriginalChunkSequenceNumber: req.ChunkSequenceNumber,
		OriginalFileType:            req.FileType,
	}

	return &processingResult, nil
}

func processStreamUnknownFileReconRequest() (*recon_responses.StreamFileChunkForReconResponse, error) {
	return nil, errors.New("Unknown File Type Supplied")
}
