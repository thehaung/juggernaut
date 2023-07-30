package httputil

import (
	"encoding/json"
	"github.com/thehaung/juggernaut/domain"
	"log"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

func WriteError(w http.ResponseWriter, errorCode int, err error) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(errorCode)

	errorResponse := domain.ErrorResponse{
		ErrorMessage: err.Error(),
		ErrorCode:    errorCode,
	}

	res, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println("httputil - WriteError - json.Marshal. Error:", err.Error())
	}

	if _, err = w.Write(res); err != nil {
		log.Println("httputil - WriteError - w.Write. Error:", err.Error())
	}
}

func WriteResponse(w http.ResponseWriter, errorCode int, response interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(errorCode)

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("httputil - WriteResponse - json.Marshal. Error:", err.Error())
	}

	if _, err = w.Write(res); err != nil {
		log.Println("httputil - WriteResponse - w.Write. Error:", err.Error())
	}
}
