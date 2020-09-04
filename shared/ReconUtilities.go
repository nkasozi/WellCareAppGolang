package shared

import (
	"encoding/json"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"log"
	"net/http"
	"strings"
)

func ToJsonString(resp interface{}) (string, error) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GenerateCORsResponse(response http.ResponseWriter, allowedMethods string) {
	response.Header().Set("Access-Control-Allow-Origin", Constants.CORS_ALLOWED_ORIGINS)
	response.Header().Set("Access-Control-Allow-Methods", allowedMethods)
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Set("Access-Control-Max-Age", "3600")
	response.WriteHeader(http.StatusNoContent)
	return
}

func GenerateMethodNotAllowedResponse(response http.ResponseWriter, method string) {
	errorMsg := method + " Is Not Supported for this API"
	response.Header().Set("Content-Type", "text/plain")
	response.Header().Set("Access-Control-Allow-Origin", Constants.CORS_ALLOWED_ORIGINS)
	response.WriteHeader(http.StatusMethodNotAllowed)
	response.Write([]byte(errorMsg))
	log.Println(errorMsg)
	return
}

func GenerateBadRequestResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", Constants.CORS_ALLOWED_ORIGINS)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
	log.Println(err.Error())
	return
}

func GetHttpScheme(request *http.Request) string {
	if strings.Contains(request.Host,"localhost") {
		return "http://"
	}

	return "https://"
}

func GenerateRedirectResponse(w http.ResponseWriter, request *http.Request, newUrl string) {
	http.Redirect(w, request, newUrl, http.StatusSeeOther)
}

func GenerateOkRequestResponse(w http.ResponseWriter, json string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", Constants.CORS_ALLOWED_ORIGINS)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func GenerateInternalServerResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", Constants.CORS_ALLOWED_ORIGINS)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Println(err.Error())
	return
}
