package webapi

import (
	_ "gitlab.com/capslock-ltd/reconciler/backend-golang/docs"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func Swagger(response http.ResponseWriter, request *http.Request) {


	switch request.Method {

	case http.MethodOptions:
		shared.GenerateCORsResponse(response, http.MethodGet)
		return

	case http.MethodGet:
		swaggerDocUrl := shared.GetHttpScheme(request) + request.Host + "/Swagger/swagger.json"
		finalRedirectUrl:= Constants.SWAGGER_EDITOR_URL + "?url=" + swaggerDocUrl
		shared.GenerateRedirectResponse(response, request,finalRedirectUrl)
		return

	default:
		shared.GenerateMethodNotAllowedResponse(response, request.Method)
		return
	}
}

func SwaggerDoc(response http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodOptions:

		//handle CORs preflight request
		shared.GenerateCORsResponse(response, http.MethodPost)
		return

	case http.MethodGet:

		//read file with json
		executablesPath, execErr := os.Executable()
		if execErr != nil {
			log.Println(execErr)
		}

		filePathWithCorrectSlashes := filepath.Dir(executablesPath) + filepath.FromSlash("/docs/swagger.json")
		dataBytes, err := ioutil.ReadFile(filePathWithCorrectSlashes)

		//error reading file
		if err != nil {
			shared.GenerateInternalServerResponse(response, err)
			return
		}

		//echo out the json
		dataString := string(dataBytes)
		shared.GenerateOkRequestResponse(response, dataString)
		return

	default:
		//oops unhandled method invoked
		shared.GenerateMethodNotAllowedResponse(response, request.Method)
		return
	}
}
