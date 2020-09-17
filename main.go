package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/webapi"
)

// @title Reconciler Backend Core API
// @version 1.0
// @description This is the core backend Reconciler API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host us-central1-reconcilercore.cloudfunctions.net
// @BasePath /
func main() {

	//setup a new router
	router := httprouter.New()

	//setup CORs globally
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shared.GenerateCORsResponse(w, "POST,GET,PUT,DELETE,OPTIONS")
		return
	})

	//handle routes
	router.POST("/GetFileUploadParameters", GetFileUploadParameters)
	router.POST("/StreamFileChunksForRecon", StreamFileChunksForRecon)
	router.GET("/Swagger/index.html", Swagger)
	router.GET("/Swagger/swagger.json", SwaggerDoc)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = Constants.DEFAULT_SERVER_PORT
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}

// GetFileUploadParameters godoc
// @Summary GetFileUploadParameters
// @Description given certain details about an incoming upload, it retrieves information necessary for successfull upload e.g batch size
// @Tags GetFileUploadParameters API
// @Accept json
// @Produce json
// @Param request body recon_requests.GetFileUploadParametersRequest true "GetFileUploadParametersRequest"
// @Success 200 {object} recon_responses.GetFileUploadParametersResponse "GetFileUploadParametersResponse"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /GetFileUploadParameters [post]
func GetFileUploadParameters(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	webapi.GetFileUploadParameters(w, r)
	return
}

// StreamFileChunksForRecon godoc
// @Summary StreamFileChunksForRecon
// @Description Receives either Source or Comparison File Chunks and routes them appropriately for Reconciliation
// @Tags StreamFileChunksForRecon API
// @Accept json
// @Produce json
// @Param request body recon_requests.StreamFileChunkForReconRequest true "StreamFileChunkForReconRequest"
// @Success 200 {object} recon_responses.StreamFileChunkForReconResponse "StreamFileChunkForReconResponse"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /StreamFileChunksForRecon [post]
func StreamFileChunksForRecon(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	webapi.StreamFileChunksForRecon(w, r)
	return
}

// GetSwaggerJson godoc
// @Summary GetSwaggerJson
// @Description returns json needed by Swagger
// @Tags Swagger APIs
// @Produce json
// @Success 200 {string} string "Json data"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /Swagger/swagger.json [get]
func SwaggerDoc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	webapi.SwaggerDoc(w, r)
	return
}

// Swagger API godoc
// @Summary Swagger API
// @Description used to access the swagger GUI
// @Tags Swagger APIs
// @Produce json
// @Success 200 {string} string "Json data"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /Swagger/index.html [get]
func Swagger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	webapi.Swagger(w, r)
	return
}
