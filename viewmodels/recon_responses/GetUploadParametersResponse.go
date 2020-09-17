package recon_responses

type GetFileUploadParametersResponse struct {
	UploadRequestId                 string
	SourceFileHash                  string
	SourceFileName                  string
	SourceFileExpectedBatchSize     int
	ComparisonFileHash              string
	ComparisonFileName              string
	ComparisonFileExpectedBatchSize int
	IsFirstTimeUploadForSourceFile  bool
	SourceFileLastRowReceived       int
	IsFirstTimeUploadForCmpFile     bool
	ComparisonFileLastRowReceived   int
}
