package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var DNSServer = "1.1.1.1:53"

func main() {
	ourDomain := ".hanami.magicmail.run."
	if os.Getenv("HANAMI_ENV") == "dev" {
		ourDomain = ".dev.hanami.run."
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		hostname := strings.Split(r.Host, ":")[0]

		// Look up cname to find redirect rule
		cname, err := net.LookupCNAME(hostname)
		if err != nil {
			w.Write([]byte("Fail dns look up. Retry"))
			return
		}

		fmt.Println("Found cname", cname)
		if !strings.HasSuffix(cname, ourDomain) {
			w.Write([]byte("Invalid request. Did you forgot to convert"))
			return
		}

		ourHostPosition := strings.Index(cname, ourDomain)

		userDomainName := cname[0:ourHostPosition]

		//originalDomain := strings.ReplaceAll(underscoreOriginalHostname, "_", ".")
		redirectTo := fmt.Sprintf("https://%s%s", userDomainName, r.URL.Path)
		fmt.Println("go to", redirectTo)

		http.Redirect(w, r, redirectTo, 302)
	})

	http.ListenAndServe(":4040", r)

}
