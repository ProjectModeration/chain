package server

import (
	"net/http"
)

func Start() {
	http.ListenAndServe(":8080", nil)
}
