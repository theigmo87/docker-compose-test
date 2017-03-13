package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var subdomains = [...]string{"taxslayer"}

// Prox struct
type Prox struct {
	// target url of reverse proxy
	target *url.URL
	// instance of Go ReverseProxy thatwill do the job for us
	proxy *httputil.ReverseProxy
}

// New small factory
func New(target string) *Prox {
	url, _ := url.Parse(target)
	// you should handle error on parsing
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func isValidSubdomain(host string) bool {
	hostParts := strings.Split(host, ".")
	log.Printf("Hostname: %s", host)
	fmt.Printf("Hostname: %s", host)
	if len(hostParts) != 3 {
		return false
	}
	for _, s := range subdomains {
		if strings.EqualFold(s, hostParts[0]) {
			return true
		}
	}
	return false

}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	if isValidSubdomain(r.Host) {
		// w.Header().Set("X-ProxyContainerId", os.Getenv("HOSTNAME"))
		p.proxy.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "http://workful.local", http.StatusSeeOther)
	}
}

func main() {
	// come constants and usage helper
	// subdomainPort := os.Getenv("SUBDOMAINSERVER_PORT")
	const (
		defaultPort        = ":8083"
		defaultPortUsage   = "default server port, ':8083'..."
		defaultTarget      = "http://subdomainserver"
		defaultTargetUsage = "default redirect url, 'http://subdomainserver'"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	url := flag.String("url", defaultTarget, defaultTargetUsage)

	flag.Parse()

	fmt.Println("server will run on : %s", *port)
	fmt.Println("redirecting to :%s", *url)

	// proxy
	proxy := New(*url)

	// server
	http.HandleFunc("/", proxy.handle)
	http.ListenAndServe(*port, nil)
}
