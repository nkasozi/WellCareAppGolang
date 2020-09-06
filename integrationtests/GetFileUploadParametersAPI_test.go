package integrationtests

import (
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon-requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/webapi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//swaggo
func TestGetFileUploadParmatersGivenValidRequestExpectSuccess(t *testing.T) {

	Convey("Check that when given a Valid GetFileUploadRequest, API Returns Success Response", t, func() {

		Convey("Given that we have setup the request", func() {

			//build the request
			bodyStruct := recon_requests.GetFileUploadParametersRequest{
				SourceFileName:            "MyTestSourceFile.csv",
				SourceFileHash:            "616161761",
				SourceFileColumnCount:     3,
				SourceFileRowCount:        100000,
				ComparisionFileName:       "MyTestComparisonFile.csv",
				ComparisonFileHash:        "7277818818",
				ComparisonFileColumnCount: 3,
				ComparisonFileRowCount:    2000000,
				ComparisonPairs: []recon_requests.ComparisonPair{{
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
				So(requestErr,ShouldEqual,nil)

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
					So(resp.StatusCode,ShouldEqual,http.StatusOK)

					//read the  body
					body, readBodyErr := ioutil.ReadAll(resp.Body)

					//oops reading body returned an error stop
					So(readBodyErr,ShouldEqual,nil)

					//decode the json response
					jsonResponse := string(body)

					//decoding failed
					So(len(jsonResponse), ShouldBeGreaterThan, 0 )
				})
			})
		})
	})
}
