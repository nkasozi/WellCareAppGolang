package integrationtests

import (
	main2 "gitlab.com/capslock-ltd/reconciler/backend-golang/main"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/viewmodels/recon-requests"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//swaggo
func TestGetFileUploadParmatersGivenValidRequestExpectSuccess(t *testing.T) {

	bodyStruct := recon_requests.GetFileUploadParametersRequest{
		SourceFileName:            "MyTestSourceFile.csv",
		SourceFileHash:            "616161761",
		SourceFileColumnCount:     3,
		SourceFileRowCount:        100000,
		ComparisionFileName:       "MyTestComparisonFile.csv",
		ComparisonFileHash:        "7277818818",
		ComparisonFileColumnCount: "3",
		ComparisonFileRowCount:    2000000,
	}

	jsonRequest, jsonErr := shared.ToJsonString(bodyStruct)

	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	request, requestErr := http.NewRequest("POST", "/GetFileUploadParameters", strings.NewReader(jsonRequest))

	if requestErr != nil {
		t.Fatal(requestErr)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(main2.GetFileUploadParameters)

	handler.ServeHTTP(recorder, request)

	resp := recorder.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Http Status Code returned is %v yet we Expected %v", resp.StatusCode, http.StatusOK)
	}

	body, readBodyErr := ioutil.ReadAll(resp.Body)

	if readBodyErr != nil {
		t.Fatal(readBodyErr)
	}

	jsonResponse := string(body)

	if len(jsonResponse) <= 0 {
		t.Errorf("Empty Response returned by Server")
	}

}
