package Helpers

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
)



func SetUpInvalidGetFileUploadParametersRequest() string {
	//build the request
	bodyStruct := recon_requests.GetFileUploadParametersRequest{}

	//turn the reuqets into json
	jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

	//oops something went wrong on turning
	So(jsonErr, ShouldEqual, nil)
	return jsonRequest
}

func AssertThatTheGetFileUploadParametersResponseWasSuccessful(resp *http.Response) *recon_responses.GetFileUploadParametersResponse {
	//status code should be success
	So(resp.StatusCode, ShouldEqual, http.StatusOK)

	//read the  body
	body, readBodyErr := ioutil.ReadAll(resp.Body)

	//oops reading body returned an error stop
	So(readBodyErr, ShouldEqual, nil)

	uploadParameters := recon_responses.GetFileUploadParametersResponse{}

	unmarshallErr := jsonHelper.Unmarshal(body, &uploadParameters)

	So(unmarshallErr, ShouldEqual, nil)

	return &uploadParameters
}

func SendTestGetFileUploadParametersRequest(jsonRequest string) *http.Response {
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
	return resp
}

func SetUpValidGetFileUploadParametersRequest() string {
	bodyStruct := recon_requests.GetFileUploadParametersRequest{
		UserId:                 shared.GenerateUniqueId("TestUser"),
		SourceFileName:         "MyTestSourceFile.csv",
		SourceFileHash:         shared.GenerateUniqueId("TestSrcHash-"),
		SourceFileRowCount:     200,
		ComparisionFileName:    "MyTestComparisonFile.csv",
		ComparisonFileHash:     shared.GenerateUniqueId("TestComparisonHash-"),
		ComparisonFileRowCount: 200,
		ComparisonPairs: []Entities.ComparisonPair{{
			ComparisonColumnIndex: 0,
			SourceColumnIndex:     1,
		}},
	}

	//turn the reuqets into json
	jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

	//oops something went wrong on turning
	So(jsonErr, ShouldEqual, nil)

	return jsonRequest
}
