package webapi

import (
	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/capslock-ltd/reconciler/backend-golang/docs"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
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
	swaggerDocUrl := "http://" + request.Host + "/Swagger/swagger.json"
	httpSwagger.Handler(httpSwagger.URL(swaggerDocUrl))
}

func SwaggerDoc(response http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodOptions:

		//handle CORs preflight request
		shared.GenerateCORsResponse(response, http.MethodPost)
		return

	case http.MethodGet:

		//read file with json
		//filePathWithCorrectSlashes := filepath.FromSlash("./docs/swagger.json")
		//dataBytes, err := ioutil.ReadFile(filePathWithCorrectSlashes)

		files:=""
		serr := filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				files+=path+","
				return nil
			})
		if serr != nil {
			log.Println(serr)
		}
		shared.GenerateOkRequestResponse(response, files)
		return
		//error reading file
		//if err != nil {
		//	shared.GenerateInternalServerResponse(response, err)
		//	return
		//}
		//
		////echo out the json
		//dataString := string(dataBytes)
		//shared.GenerateOkRequestResponse(response, dataString)
		//return

	default:
		//oops unhandled method invoked
		shared.GenerateMethodNotAllowedResponse(response, request.Method)
		return
	}
}
