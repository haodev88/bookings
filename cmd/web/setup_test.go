package main

import (
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M)  {
	os.Exit(m.Run())
}

type myHandler struct {

}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Run serveHTTP")
}
