package integrationtests

import (
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/integrationtests/Helpers"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"testing"
)

func TestThatGivenValidStreamReconComparisonFileRequestExpectSuccess(t *testing.T) {

	Convey("Check that when given a Valid StreamFileChunkForReconRequest for a Source File, API Returns Success Response", t, func() {

		Convey("Given that we have setup a Valid GetFileUploadRequest request", func() {

			//build the request
			jsonRequest := Helpers.SetUpValidGetFileUploadParametersRequest()

			Convey("And we send the request", func() {

				resp := Helpers.SendTestGetFileUploadParametersRequest(jsonRequest)

				Convey("Then the response should be success", func() {

					uploadParameters := Helpers.AssertThatTheGetFileUploadParametersResponseWasSuccessful(resp)

					//We can now proceed and test the recon API
					Convey("And Given that we have a Valid StreamFileChunkForRecon request", func() {

						jsonRequest := Helpers.SetUpValidStreamFileChunkForReconRequest(uploadParameters, UploadFileTypes.ComparisonFile)

						Convey("And we send the request", func() {

							resp := Helpers.SendTestStreamFileChunkForReconRequest(jsonRequest)

							Convey("Then the response should be success", func() {

								Helpers.AssertThatTheStreamFileChunkForReconResponseWasSuccessful(resp)

							})
						})
					})

				})
			})
		})

	})
}

func TestThatGivenValidStreamReconSrcFileRequestExpectSuccess(t *testing.T) {

	Convey("Check that when given a Valid StreamFileChunkForReconRequest for a Source File, API Returns Success Response", t, func() {

		Convey("Given that we have setup a Valid GetFileUploadRequest request", func() {

			//build the request
			jsonRequest := Helpers.SetUpValidGetFileUploadParametersRequest()

			Convey("And we send the request", func() {

				resp := Helpers.SendTestGetFileUploadParametersRequest(jsonRequest)

				Convey("Then the response should be success", func() {

					uploadParameters := Helpers.AssertThatTheGetFileUploadParametersResponseWasSuccessful(resp)

					//We can now proceed and test the recon API
					Convey("And Given that we have a Valid StreamFileChunkForRecon request", func() {

						jsonRequest := Helpers.SetUpValidStreamFileChunkForReconRequest(uploadParameters, UploadFileTypes.SrcFile)

						Convey("And we send the request", func() {

							resp := Helpers.SendTestStreamFileChunkForReconRequest(jsonRequest)

							Convey("Then the response should be success", func() {

								Helpers.AssertThatTheStreamFileChunkForReconResponseWasSuccessful(resp)

							})
						})
					})

				})
			})
		})

	})
}

func TestThatGivenInvalidStreamReconFileRequestExpectSuccess(t *testing.T) {

	Convey("Check that when given a Valid StreamFileChunkForReconRequest for a Source File, API Returns Success Response", t, func() {

		//We can now proceed and test the recon API
		Convey("And Given that we have an Invalid StreamFileChunkForRecon request", func() {

			jsonRequest := Helpers.SetUpInvalidStreamFileChunkForReconRequest()

			Convey("And we send the request", func() {

				resp := Helpers.SendTestStreamFileChunkForReconRequest(jsonRequest)

				Convey("Then the response should be success", func() {

					Helpers.AssertThatResponseWasNotSuccessful(resp)

				})
			})
		})

	})
}
