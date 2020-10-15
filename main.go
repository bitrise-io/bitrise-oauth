package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/bitrise-io/bitrise-oauth/client"
	"github.com/bitrise-io/bitrise-oauth/service"
)

func testCall() {
	time.Sleep(time.Second * 3)
	fmt.Println("# get")

	authProvider := client.NewWithSecret("tomi-test", "37dd1bc5-50bb-4674-a5fa-2cec87037e52")

	resp, err := authProvider.ManagedHTTPClient().Get("http://localhost:8080/test")
	if err != nil {
		panic(err)
	}

	rb, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("resp:\n" + string(rb))
}

func handler(w http.ResponseWriter, r *http.Request) {
	rb, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rb))
}

func main() {
	go testCall()

	mux := http.NewServeMux()

	kti := service.NewValidator()

	mux.Handle("/test", kti.Middleware(http.HandlerFunc(handler)))

	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
