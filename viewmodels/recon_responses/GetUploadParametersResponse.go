package recon_responses


type GetFileUploadParametersResponse struct {
	SourceFileHash              string
	SourceFileName              string
	SourceFileExpectedBatchSize int

	ComparisonFileHash              string
	ComparisonFileName              string
	ComparisonFileExpectedBatchSize int
}


