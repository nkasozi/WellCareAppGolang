package integrationtests

import (
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/integrationtests/Helpers"

	"testing"
)

//convey tests
func TestGetFileUploadParametersGivenValidRequestExpectSuccess(t *testing.T) {

	Convey("Check that when given a Valid GetFileUploadRequest, API Returns Success Response", t, func() {

		Helpers.PointPubSubServiceProviderToAnInMemoryInstance()

		Convey("Given that we have setup a Valid GetFileUploadRequest request", func() {

			//build the request
			jsonRequest := Helpers.SetUpValidGetFileUploadParametersRequest()

			Convey("And we send the request", func() {

				resp :=  Helpers.SendTestGetFileUploadParametersRequest(jsonRequest)

				Convey("Then the response should be success", func() {

					Helpers.AssertThatTheGetFileUploadParametersResponseWasSuccessful(resp)

				})
			})
		})
	})
}



//convy tsts
func TestGetFileUploadParametersGivenInvalidRequestExpectFailure(t *testing.T) {

	Convey("Check that when given a invalid GetFileUploadRequest, API Returns Failure Response", t, func() {

		Helpers.PointPubSubServiceProviderToAnInMemoryInstance()

		Convey("Given that we have setup the request", func() {

			jsonRequest :=  Helpers.SetUpInvalidGetFileUploadParametersRequest()

			Convey("And we send the request", func() {

				resp :=  Helpers.SendTestGetFileUploadParametersRequest(jsonRequest)

				Convey("Then the response should NOT be success", func() {

					Helpers.AssertThatResponseWasNotSuccessful(resp)

				})
			})
		})
	})
}


