package Entities

import (
	"time"
)

type FilesUploadedParameters struct{

	//id of user uploading file
	UserId string

	//src file meta data
	SourceFileName string
	SourceFileHash string
	SourceFileRowCount int
	SourceFileExpectedBatchSize int
	SourceFileLastRowReceived int

	//comparison file meta data
	ComparisionFileName string
	ComparisonFileHash string
	ComparisonFileRowCount int
	ComparisonFileExpectedBatchSize int
	ComparisonFileLastRowReceived int

	//comparison data
	ComparisonPairs []ComparisonPair

	//pubsub fields
	SourceFileTopics []string
	ComparisonFileTopic string

	//db model fields
	Id string
	DateCreated time.Time
	DateModified time.Time
}

type ComparisonPair struct {
	SourceColumnIndex int `validate:"gte=0"`
	ComparisonColumnIndex int `validate:"gte=0"`
}

func (fu FilesUploadedParameters) GetEntityId() string {
	return fu.Id
}

func (fu FilesUploadedParameters) SetEntityId(Id string) EntityInterface {
	fu.Id = Id
	return  fu
}

func (fu FilesUploadedParameters) GetDateCreated() string {
	return "2020-08-20"
}

func (fu FilesUploadedParameters) GetDateModified() string {
	return "2020-08-20"
}
