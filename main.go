package main

import (
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/webapi"
	"log"
	"net/http"
	"os"
)

// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	//handle routes
	http.HandleFunc("/GetFileUploadParameters", webapi.GetFileUploadParametersNew)
	http.HandleFunc("/Swagger", webapi.Swagger)
	http.HandleFunc("/Swagger/swagger.json", webapi.SwaggerDoc)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = Constants.DEFAULT_SERVER_PORT
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

