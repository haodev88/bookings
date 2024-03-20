package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNosurf(t *testing.T) {
	var mH myHandler
	h:= Nosurf(&mH)
	switch v:= h.(type) {
	case http.Handler:
		// do something
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myHand myHandler
	mySessionLoad:= SessionLoad(&myHand)
	switch v:= mySessionLoad.(type) {
		case http.Handler:
		// do something
		default:
		t.Error(fmt.Sprintf("Type session load not http handler but is %T", v))
	}
}