package Entities

import (
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"time"
)

type StreamedFileChunk struct {
	//db model fields
	Id                  string
	DateCreated         time.Time
	DateModified        time.Time
	FileType            UploadFileTypes.UploadFileType
	UploadRequestId     string
	ChunkSequenceNumber int
	Records             []string
	IsEOF               bool
}

func (sf StreamedFileChunk) GetEntityId() string {
	return sf.Id
}

func (sf StreamedFileChunk) SetEntityId(Id string) EntityInterface {
	sf.Id = Id
	return sf
}

func (sf StreamedFileChunk) GetDateCreated() string {
	return "2020-08-20"
}

func (StreamedFileChunk) GetDateModified() string {
	return "2020-08-20"
}
