package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

type clientInfo struct {
	RemoteAddr string
}

func writeClientTemplate(w io.Writer, ci *clientInfo, tpl *template.Template) {
	tpl.Execute(w, ci)
}

func writeClientText(w io.Writer, ci *clientInfo) {
	fmt.Fprintf(w, "remote address: %s", ci.RemoteAddr)
}

func getClientInfo(r *http.Request) *clientInfo {
	return &clientInfo{
		RemoteAddr: r.RemoteAddr,
	}
}

func handleRquestInfo(w http.ResponseWriter, r *http.Request) {
	ci := getClientInfo(r)
	writeClientText(w, ci)
}

func main() {
	// Determine the port.
	var la string
	var host string
	var err error
	if host, err = os.Hostname(); err != nil {
		os.Exit(1)
	}
	la = ":8080"
	if p := os.Getenv("PORT"); len(p) > 0 {
		la = ":" + p
	}

	http.HandleFunc("/", handleRquestInfo)
	fmt.Printf("Listening on http://%s%s ...", host, la)
	http.ListenAndServe(la, http.DefaultServeMux)
}
