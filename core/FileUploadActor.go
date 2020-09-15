package core

import (
	"cloud.google.com/go/pubsub"
	"errors"
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core/gcloud"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/TopicSubscriberTypes"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"math"
	"strconv"
	"sync"
	"time"
)

var redisDataStore = datastore.RedisDataStoreFilesUplaodedParameters{}

func createTopicSubscribers(subscribers []Entities.TopicSubscriber) (ResultsErr error) {

	//create wait group
	var waitGroup sync.WaitGroup

	//create a channel we will use to publish results
	resultsChannel := make(chan error)

	//maek sure we always close the channel at the end
	defer close(resultsChannel)

	//loop thru all subs and create them in cloud pub sub
	// in parallel
	for _, subscriber := range subscribers {

		//increment the wait group counter by one
		waitGroup.Add(1)

		//spawn the wait group
		go func() {

			//create the subscriber
			err := createSubscriber(subscriber)

			//publish the error result
			//on the channel
			resultsChannel <- err

			//make sure we always notify the wait group
			//when we are done, no matter what
			defer waitGroup.Done()
		}()
	}

	//go through all the errors returned
	for err := range resultsChannel {

		//oops some go routine
		//failed
		if err != nil {

			//take note of any
			//serious errors
			ResultsErr = err

		}
	}

	//wait for all go routines finish
	waitGroup.Wait()

	//if one of the go routines didnt finish successfully
	//this wont be nil
	return ResultsErr
}

func createSubscriber(subscriber Entities.TopicSubscriber) error {
	if err := subscriber.IsValid(); err != nil {
		return err
	}

	switch subscriber.SubscriberType {
	case TopicSubscriberTypes.PullSubscriber:
		_, err := gcloud.CreatePullSubscriberForTopicOnCloudPubSub(Constants.GOOGLE_CLOUD_PROJECT_ID, subscriber.Name, subscriber.TopicName)
		return err

	case TopicSubscriberTypes.PushSubscriber:
		_, err := gcloud.CreatePushSubscriberForTopicOnCloudPubSub(Constants.GOOGLE_CLOUD_PROJECT_ID, subscriber.Name, subscriber.TopicName, subscriber.NotificationUrl)
		return err

	default:
		return fmt.Errorf("Unknown Subscriber Type Supplied")
	}
}

//determines the batch for the file upload
func determineBatchSizeForFile(req recon_requests.GetFileUploadParametersRequest, fileType UploadFileTypes.UploadFileType) (batchSize int, numOfBatches int) {
	switch fileType {
	case UploadFileTypes.ComparisonFile:
		if req.ComparisonFileRowCount < 200 {
			batchSize = req.ComparisonFileRowCount
			numOfBatches = 1
		} else {
			batchSize = 200
			numOfBatches = int(math.Ceil(float64(req.ComparisonFileRowCount) / float64(200)))
		}
		break
	case UploadFileTypes.SrcFile:
		if req.SourceFileRowCount < 200 {
			batchSize = req.SourceFileRowCount
			numOfBatches = 1
		} else {
			batchSize = 200
			numOfBatches = int(math.Ceil(float64(req.ComparisonFileRowCount) / float64(200)))
		}
		break

	}
	return batchSize, numOfBatches
}

//saves an entity in a datastore
func saveUploadRequest(req Entities.FilesUploadedParameters, isCreateRequest bool) (string, error) {

	//create request
	if isCreateRequest || req.Id == "0" {
		return redisDataStore.Add(req)
	}

	//update request
	return redisDataStore.Update(req)
}

func ProcessGetUploadParametersRequest(req recon_requests.GetFileUploadParametersRequest) (processingResult *recon_responses.GetFileUploadParametersResponse, err error) {

	defer func() (*recon_responses.GetFileUploadParametersResponse, error) {

		//catch any exceptions and return a clean error
		if r := recover(); r != nil {

			err = errors.New(fmt.Sprint("Error:", r))
			return nil, err
		}

		//success response return as is
		return processingResult, err
	}()

	//check if we have ever recieved these files before
	response, checkErr := checkIfIsRepeatUploadAttempt(req)

	//yes its a repeat attempt
	if checkErr == nil {
		return response, nil
	}

	//process the request
	return processRequestToUploadNewFile(req)
}

func checkIfIsRepeatUploadAttempt(req recon_requests.GetFileUploadParametersRequest) (*recon_responses.GetFileUploadParametersResponse, error) {

	//check if we have recived this item before
	fileId := generateFileUploadID(req)
	existingFileUploadAttempt, getErr := redisDataStore.GetById(fileId)

	//error in checking for item in redis
	if getErr != nil {
		return nil, getErr
	}

	//means we have no existing attempt
	if existingFileUploadAttempt.Id == "" {
		return nil, errors.New("No previous Attempt has been logged for File with ID [" + fileId + "]")
	}

	//build up response
	processingResult := recon_responses.GetFileUploadParametersResponse{}
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
	return &processingResult, nil
}

func processRequestToUploadNewFile(req recon_requests.GetFileUploadParametersRequest) (*recon_responses.GetFileUploadParametersResponse, error) {

	//determine the batch size for each upload request
	srcFileBatchSize, srcFileNumOfBatches := determineBatchSizeForFile(req, UploadFileTypes.SrcFile)
	comparisonFileBatchSize, _ := determineBatchSizeForFile(req, UploadFileTypes.ComparisonFile)

	//create source file topics in Queue Service
	srcFileTopicNames := GenerateNamesForSourceFileTopics(req, srcFileNumOfBatches)
	topicSubscribers := getSourceFileTopicSubscribers(req, srcFileTopicNames)
	_, err := createSourceFileTopicOnCloudPubSub(Constants.GOOGLE_CLOUD_PROJECT_ID, srcFileTopicNames)

	if err != nil {
		return nil, err
	}

	//create comparison file topics in Queue Service
	comparisonFileTopicName := GenerateNameForComparisonFileTopic(req.UserId, req.ComparisonFileHash)
	topicSubscribers = getComparisonFileSubscribers(topicSubscribers)
	_, err = gcloud.CreateTopicOnCloudPubSub(Constants.GOOGLE_CLOUD_PROJECT_ID, comparisonFileTopicName)

	if err != nil {
		return nil, err
	}

	//create subscribers for the topics that were created in the queue service
	err = createTopicSubscribers(topicSubscribers)

	//fialed to create all the subscribers
	if err != nil {
		return nil, err
	}

	//map to datatstore struct to save
	filesUploadedParameters := Entities.FilesUploadedParameters{
		UserId:                          req.UserId,
		SourceFileName:                  req.SourceFileName,
		SourceFileHash:                  req.SourceFileHash,
		SourceFileRowCount:              req.SourceFileRowCount,
		SourceFileExpectedBatchSize:     srcFileBatchSize,
		SourceFileLastRowReceived:       0,
		SourceFileTopics:                srcFileTopicNames,
		ComparisionFileName:             req.ComparisionFileName,
		ComparisonFileHash:              req.ComparisonFileHash,
		ComparisonFileRowCount:          req.ComparisonFileRowCount,
		ComparisonFileExpectedBatchSize: comparisonFileBatchSize,
		ComparisonFileLastRowReceived:   0,
		ComparisonFileTopic:             comparisonFileTopicName,
		ComparisonPairs:                 req.ComparisonPairs,
		Id:                              generateFileUploadID(req),
		DateCreated:                     time.Now(),
		DateModified:                    time.Now(),
	}

	//save the upload request meta data in redis
	uploadRequestId, saveErr := saveUploadRequest(filesUploadedParameters, true)

	//failed to save..stop here
	if saveErr != nil {
		return nil, saveErr
	}

	//build up response
	processingResult := recon_responses.GetFileUploadParametersResponse{}
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
	return &processingResult, nil
}

func getComparisonFileSubscribers(topicSubscribers []Entities.TopicSubscriber) []Entities.TopicSubscriber {

	for _, name := range topicSubscribers {

		topicSubscriber := Entities.TopicSubscriber{
			Name:            name.Name,
			SubscriberType:  TopicSubscriberTypes.PullSubscriber,
			NotificationUrl: "https://badadc09b6648dee9da7515d55e0ec68.m.pipedream.net",
			TopicName:       name.TopicName,
		}

		topicSubscribers = append(topicSubscribers, topicSubscriber)
	}

	return topicSubscribers
}

func getSourceFileTopicSubscribers(req recon_requests.GetFileUploadParametersRequest, srcFileTopicNames []string) []Entities.TopicSubscriber {

	var topicSubscribers []Entities.TopicSubscriber

	for i, name := range srcFileTopicNames {

		topicSubscriber := Entities.TopicSubscriber{
			Name:            GenerateNameForSourceFileTopic(req.UserId, req.SourceFileHash, i),
			SubscriberType:  TopicSubscriberTypes.PushSubscriber,
			NotificationUrl: "https://badadc09b6648dee9da7515d55e0ec68.m.pipedream.net",
			TopicName:       name,
		}

		topicSubscribers = append(topicSubscribers, topicSubscriber)
	}
	
	return topicSubscribers
}

func createSourceFileTopicOnCloudPubSub(projectId string, names []string) ([]pubsub.Topic, error) {

	var topicPtrs []pubsub.Topic

	for _, name := range names {

		topicPtr, err := gcloud.CreateTopicOnCloudPubSub(projectId, name)

		if err != nil {
			return nil, err
		}

		topicPtrs = append(topicPtrs, *topicPtr)
	}

	return topicPtrs, nil
}

func GenerateNamesForSourceFileTopics(req recon_requests.GetFileUploadParametersRequest, NumberOfNames int) []string {

	var names []string

	for i := 0; i < NumberOfNames; i++ {
		name := GenerateNameForSourceFileTopic(req.UserId, req.SourceFileHash, i)
		names = append(names, name)
	}

	return names
}

func GenerateNameForSourceFileTopic(UserId, Hash string, i int) string {
	return "SRC-" + UserId + "-" + shared.GenerateUniqueId(Hash) + "-" + strconv.Itoa(i)
}

func GenerateNameForComparisonFileTopic(UserId, Hash string) string {
	return "CMP-" + UserId + "-" + shared.GenerateUniqueId(Hash)
}

func generateFileUploadID(req recon_requests.GetFileUploadParametersRequest) string {
	return req.UserId + "-" + req.SourceFileHash + "-" + req.ComparisonFileHash
}
