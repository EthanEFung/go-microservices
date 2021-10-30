package handlers

import (
	"log"
	"net/http"
)

type OK struct {
	l *log.Logger
}

func NewOK(l *log.Logger) *OK {
	return &OK{l}
}

func (ok *OK) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
}
