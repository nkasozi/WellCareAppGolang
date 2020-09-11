package integrationtests

import (
	jsonHelper "encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/webapi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//convy tsts
func TestGetFileUploadParametersGivenValidRequestExpectSuccess(t *testing.T) {

	Convey("Check that when given a Valid GetFileUploadRequest, API Returns Success Response", t, func() {

		Convey("Given that we have setup the request", func() {

			//build the request
			bodyStruct := recon_requests.GetFileUploadParametersRequest{
				UserId:                 "TestUser67338738",
				SourceFileName:         "MyTestSourceFile.csv",
				SourceFileHash:         shared.GenerateUniqueId("TestSrcHash-"),
				SourceFileRowCount:     100000,
				ComparisionFileName:    "MyTestComparisonFile.csv",
				ComparisonFileHash:     shared.GenerateUniqueId("TestComparisonHash-"),
				ComparisonFileRowCount: 2000000,
				ComparisonPairs: []Entities.ComparisonPair{{
					ComparisonColumnIndex: 0,
					SourceColumnIndex:     1,
				}},
			}

			//turn the reuqets into json
			jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

			//oops something went wrong on turning
			So(jsonErr, ShouldEqual, nil)

			Convey("And we send the request", func() {

				//send the request across
				request, requestErr := http.NewRequest("POST", "/GetFileUploadParameters", strings.NewReader(jsonRequest))

				//sending didnt finish well
				So(requestErr, ShouldEqual, nil)

				//record the result
				recorder := httptest.NewRecorder()

				//set up the server
				handler := http.HandlerFunc(webapi.GetFileUploadParameters)

				//boot up the server
				handler.ServeHTTP(recorder, request)

				//recieve resp
				resp := recorder.Result()

				Convey("Then the response should be success", func() {
					//status code should be success
					So(resp.StatusCode, ShouldEqual, http.StatusOK)

					//read the  body
					body, readBodyErr := ioutil.ReadAll(resp.Body)

					//oops reading body returned an error stop
					So(readBodyErr, ShouldEqual, nil)

					unmarshallErr := jsonHelper.Unmarshal(body, &recon_responses.GetFileUploadParametersResponse{})

					So(unmarshallErr, ShouldEqual, nil)
				})
			})
		})
	})
}

//convy tsts
func TestGetFileUploadParametersGivenInvalidRequestExpectFailure(t *testing.T) {

	Convey("Check that when given a invalid GetFileUploadRequest, API Returns Failure Response", t, func() {

		Convey("Given that we have setup the request", func() {

			//build the request
			bodyStruct := recon_requests.GetFileUploadParametersRequest{}

			//turn the reuqets into json
			jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

			//oops something went wrong on turning
			So(jsonErr, ShouldEqual, nil)

			Convey("And we send the request", func() {

				//send the request across
				request, requestErr := http.NewRequest("POST", "/GetFileUploadParameters", strings.NewReader(jsonRequest))

				//sending didnt finish well
				So(requestErr, ShouldEqual, nil)

				//record the result
				recorder := httptest.NewRecorder()

				//set up the server
				handler := http.HandlerFunc(webapi.GetFileUploadParameters)

				//boot up the server
				handler.ServeHTTP(recorder, request)

				//recieve resp
				resp := recorder.Result()

				Convey("Then the response should NOT be success", func() {
					//status code should be success
					So(resp.StatusCode, ShouldNotEqual, http.StatusOK)
				})
			})
		})
	})
}
