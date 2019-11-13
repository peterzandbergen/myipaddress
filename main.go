package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"
)

type clientInfo struct {
	Port       string
	RequestURI string
	RemoteAddr string
	UserAgent  string
	Header     map[string][]string
}

func writeClientTemplate(w io.Writer, ci *clientInfo, tpl *template.Template) {
	tpl.Execute(w, ci)
}

func sortHeader(header map[string][]string) []string {
	var res = make([]string, 0, len(header))
	for k := range header {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func writeClientText(w io.Writer, ci *clientInfo) {
	fmt.Fprintf(w, "instance: %s\n", instanceName)
	fmt.Fprintf(w, "request uri: %s\n", ci.RequestURI)
	fmt.Fprintf(w, "listen port: %s\n", ci.Port)
	fmt.Fprintf(w, "server hostname: %s\n", hostname)
	fmt.Fprintf(w, "remote address: %s\n", ci.RemoteAddr)
	fmt.Fprintf(w, "UserAgent: %s\n", ci.UserAgent)
	fmt.Fprint(w, "\nHeaders:\n")
	keys := sortHeader(ci.Header)
	for _, k := range keys {
		fmt.Fprintf(w, "    %s: ", k)
		var first = true
		for _, vv := range ci.Header[k] {
			fmt.Fprintf(w, "%s", vv)
			if !first {
				fmt.Fprint(w, ", ")
			}
			first = false
		}
		fmt.Fprint(w, "\n")
	}
}

func getClientInfo(r *http.Request) *clientInfo {

	return &clientInfo{
		RequestURI: r.RequestURI,
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Header:     r.Header,
		Port:       port,
	}
}

func handleRequestInfo(w http.ResponseWriter, r *http.Request) {
	ci := getClientInfo(r)
	w.Header().Add("Content-Type", "text/plain")
	writeClientText(w, ci)
}

func logRequestInfo(w http.ResponseWriter, r *http.Request) {
	ci := getClientInfo(r)
	writeClientText(os.Stdout, ci)
	os.Stdout.Write([]byte("\n============================\n"))
}

var instanceName string = "instance"
var port string = "8080"
var hostname string

func main() {
	// Determine the port.
	var err error
	if hostname, err = os.Hostname(); err != nil {
		os.Exit(1)
	}
	if p := os.Getenv("PORT"); len(p) > 0 {
		port = p
	}
	la := ":" + port
	if i := os.Getenv("NAME"); len(i) > 0 {
		instanceName = i
	}

	http.HandleFunc("/", handleRequestInfo)
	http.HandleFunc("/nifi/", logRequestInfo)
	fmt.Printf("Listening on http://%s%s ...\n", hostname, la)
	log.Fatal(http.ListenAndServe(la, http.DefaultServeMux).Error())
}
