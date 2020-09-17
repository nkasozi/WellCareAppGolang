package Helpers

import (
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/core/pubsubserviceproviders/inmemory"
	"net/http"
)
import . "github.com/smartystreets/goconvey/convey"

func AssertThatResponseWasNotSuccessful(resp *http.Response) {
	//status code should be success
	So(resp.StatusCode, ShouldNotEqual, http.StatusOK)
}

func PointPubSubServiceProviderToAnInMemoryInstance() error{
	core.PubSubServiceProvider = inmemory.NewPubSubClient()
	return nil
}
