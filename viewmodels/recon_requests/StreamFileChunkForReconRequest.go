package recon_requests

import "gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"

type StreamFileChunkForReconRequest struct {
	FileType            UploadFileTypes.UploadFileType `validate:"required"`
	UploadRequestId     string                         `validate:"required"`
	ChunkSequenceNumber int                            `validate:"required,gte=1"`
	Records             []string                       `validate:"required"`
	IsEOF               bool                           `validate:"required"`
}

func (req StreamFileChunkForReconRequest) IsValid() error {
	return nil
}
