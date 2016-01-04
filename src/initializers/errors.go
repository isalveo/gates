package initializers

import (
	"errors"
	"net/http"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

var StatusNotFound = errors.New("Not Found")

func Error(w http.ResponseWriter, err error, status int) {
	http.Error(w, ResponseErrorMsg(err, status), status)
}

func ResponseErrorMsg(err error, status int) string {
	switch status {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 406:
		return "Not Acceptable"
	case 404:
		return "Not Found"
	case 415:
		return "Unsupported Media Type"
	case 422:
		return "Unprocessable Entity"
	case 500:
		return "Internal Server Error"
	}

	return "Unknown Error Occured"
}
