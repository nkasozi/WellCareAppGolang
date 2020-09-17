package Helpers

import (
	jsonHelper "encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/UploadFileTypes"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_requests"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon_responses"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/webapi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func AssertThatTheStreamFileChunkForReconResponseWasSuccessful(resp *http.Response) *recon_responses.StreamFileChunkForReconResponse{

	//http status code should be success
	So(resp.StatusCode, ShouldEqual, http.StatusOK)

	//read the  body
	body, readBodyErr := ioutil.ReadAll(resp.Body)

	//oops reading body returned an error stop
	So(readBodyErr, ShouldEqual, nil)

	streamFileResp:=recon_responses.StreamFileChunkForReconResponse{}

	unmarshallErr := jsonHelper.Unmarshal(body, &streamFileResp)

	So(unmarshallErr, ShouldEqual, nil)

	return &streamFileResp
}

func SendTestStreamFileChunkForReconRequest(jsonRequest string) *http.Response {
	//send the request across
	request, requestErr := http.NewRequest("POST", "/StreamFileChunksForRecon", strings.NewReader(jsonRequest))

	//sending didnt finish well
	So(requestErr, ShouldEqual, nil)

	//record the result
	recorder := httptest.NewRecorder()

	//set up the server
	handler := http.HandlerFunc(webapi.StreamFileChunksForRecon)

	//boot up the server
	handler.ServeHTTP(recorder, request)

	//recieve resp
	resp := recorder.Result()
	return resp
}

func SetUpValidStreamFileChunkForReconRequest(uploadParameters *recon_responses.GetFileUploadParametersResponse, fileType UploadFileTypes.UploadFileType) string {

	//build the request
	bodyStruct := recon_requests.StreamFileChunkForReconRequest{
		FileType:            fileType,
		UploadRequestId:     uploadParameters.UploadRequestId,
		ChunkSequenceNumber: 1,
		IsEOF: true,
		Records:             []string{"23,25000,17/02/2020,Success"},
	}

	//turn the reuqets into json
	jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

	//oops something went wrong on turning
	So(jsonErr, ShouldEqual, nil)
	return jsonRequest
}

func SetUpInvalidStreamFileChunkForReconRequest() string {

	//build the request
	bodyStruct := recon_requests.StreamFileChunkForReconRequest{
	}

	//turn the reuqets into json
	jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

	//oops something went wrong on turning
	So(jsonErr, ShouldEqual, nil)

	return jsonRequest
}
