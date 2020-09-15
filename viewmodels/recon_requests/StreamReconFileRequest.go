package recon_requests

import "gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"

type StreamReconRequest struct {
	FileType UploadFileTypes.UploadFileType
	UploadRequestId string
	ChunkSequenceNumber int
	Records []string
}

func (req StreamReconRequest) IsValid() error {
	return nil
}
