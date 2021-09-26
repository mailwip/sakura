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
	isDev := os.Getenv("HANAMI_ENV") == "dev"

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		hostname := strings.Split(r.Host, ":")[0]

		// Look up txt find redirect rule
		txt, err := net.LookupTXT(fmt.Sprintf("hanami-forward.%s", hostname))
		if err != nil {
			w.Write([]byte(fmt.Sprintf(instruction, hostname)))
			return
		}

		if len(txt) == 0 {
			w.Write([]byte(fmt.Sprintf(instruction, hostname)))
			return
		}

		redirectTo := txt[0]

		if !strings.HasPrefix(redirectTo, "http") {
			redirectTo = "http://" + redirectTo
		}

		if r.URL.Path != "/" {
			redirectTo = redirectTo + r.URL.Path
		}
		fmt.Printf("Found redirect rule for %s to %s at %s\n", hostname, redirectTo, r.URL.Path)

		if r.URL.RawQuery != "" {
			redirectTo = redirectTo + "?" + r.URL.RawQuery
		}

		http.Redirect(w, r, redirectTo, 302)
	})

	if isDev {
		http.ListenAndServe(":4040", r)
	} else {
		http.ListenAndServe(":4080", r)
	}

}

const instruction = `To rediret your URL with hanami.run forwarding service, you will need to create a TXT record on your domain and point to the redirectURL

TYPE: TXT
Hostname: hanami-forward.%s
Value: [the-url-to-redirect-to]

Docs: https://hanami.run/docs/url-forwarding
`
