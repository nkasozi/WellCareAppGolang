package recon_responses

type GetFileUploadParametersResponse struct {
	//unique Id tagged to this whole uplaod
	UploadRequestId string

	//SRC File meta data
	SourceFileHash              string
	SourceFileName              string
	SourceFileExpectedBatchSize int

	//Comparison File meta
	ComparisonFileHash              string
	ComparisonFileName              string
	ComparisonFileExpectedBatchSize int

	//flags indicating whether the SRC file
	//is new to us..if not where did we stop
	SourceFileIsFirstTimeUpload bool
	SourceFileLastRowReceived   int

	//flags indicating whether the Comparison file
	//is new to us..if not where did we stop
	ComparisonFileIsFirstTimeUpload bool
	ComparisonFileLastRowReceived   int
}
