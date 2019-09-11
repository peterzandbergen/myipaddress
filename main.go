package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type clientInfo struct {
	RemoteAddr string
	UserAgent  string
	Header     map[string][]string
}

func writeClientTemplate(w io.Writer, ci *clientInfo, tpl *template.Template) {
	tpl.Execute(w, ci)
}

func writeClientText(w io.Writer, ci *clientInfo) {
	fmt.Fprintf(w, "remote address: %s\n", ci.RemoteAddr)
	fmt.Fprintf(w, "UserAgent: %s\n", ci.UserAgent)
	for k, v := range ci.Header {
		fmt.Fprintf(w, "%s:\n", k)
		for _, vv := range v {
			fmt.Fprintf(w, "    %s\n", vv)
		}
	}
}

func getClientInfo(r *http.Request) *clientInfo {
	return &clientInfo{
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Header:     r.Header,
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
	log.Fatal(http.ListenAndServe(la, http.DefaultServeMux).Error())
}
