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
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Header:     r.Header,
	}
}

func handleRquestInfo(w http.ResponseWriter, r *http.Request) {
	ci := getClientInfo(r)
	w.Header().Add("Content-Type", "text/plain")
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
