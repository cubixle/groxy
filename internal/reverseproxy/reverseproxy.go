package reverseproxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/cubixle/groxy/internal/metrics"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Proxy struct {
	endpoints map[string]Endpoint
	logger    *zap.Logger
	client    *http.Client
}

type File struct {
	Debug     bool       `yaml:"debug"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Addr       string `yaml:"addr"`
	RemoteAddr string `yaml:"remote_addr"`
}

func NewReverseProxy(endpoints []Endpoint, logger *zap.Logger) *Proxy {
	epmap := map[string]Endpoint{}
	for _, e := range endpoints {
		epmap[e.Addr] = e
	}

	httpClient := &http.Client{}

	return &Proxy{
		endpoints: epmap,
		logger:    logger,
		client:    httpClient,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	ep, ok := p.endpoints[host]
	if !ok {
		p.logger.Error("endpoint not found",
			zap.String("host", host),
		)
		return
	}

	newURL, err := buildReqURL(ep.RemoteAddr, r.URL)
	if err != nil {
		logrus.Error(err)
		return
	}

	p.logger.Debug("building request",
		zap.String("url", newURL.String()),
	)

	req, err := http.NewRequest(r.Method, newURL.String(), r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	req.Host = newURL.Host
	req.Header.Set("Origin", newURL.String())

	copyHeader(req.Header, r.Header)

	delHopHeaders(req.Header)

	clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		appendHostToXForwardHeader(req.Header, clientIP)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}

	if err == nil {
		defer resp.Body.Close()
	}

	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.Header != nil {
		copyHeader(w.Header(), resp.Header)
	}

	mimetype := ContentType(resp)
	w.Header().Set("Content-Type", mimetype)

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)

	metrics.RequestInc(resp.StatusCode, req.Method)
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
	"Access-Control-Allow-Origin",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		k = http.CanonicalHeaderKey(k)
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}
