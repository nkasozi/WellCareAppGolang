package Helpers

import "net/http"
import . "github.com/smartystreets/goconvey/convey"

func AssertThatResponseWasNotSuccessful(resp *http.Response) {
	//status code should be success
	So(resp.StatusCode, ShouldNotEqual, http.StatusOK)
}
