package recon_responses

import "gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"

type StreamFileChunkForReconResponse struct {
	MessageId string
	Status string
	OriginalChunkSequenceNumber int
	OriginalFileType UploadFileTypes.UploadFileType
	//static channel used by all instances
}
