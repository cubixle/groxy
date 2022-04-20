package reverseproxy

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

func ContentType(rsp *http.Response) string {
	mimetype, _, _ := mime.ParseMediaType(rsp.Header.Get("Content-Type"))
	if mimetype == "" {
		b := bufio.NewReader(rsp.Body)
		rsp.Body = ioutil.NopCloser(b)
		mimetype = peekContentType(b)
	}

	return mimetype
}

// peekContentType peeks at the first 512 bytes of p, and attempts to detect
// the content type.  Returns empty string if error occurs.
func peekContentType(p *bufio.Reader) string {
	byt, err := p.Peek(512)
	if err != nil && !errors.Is(err, bufio.ErrBufferFull) && !errors.Is(err, io.EOF) {
		return ""
	}
	return http.DetectContentType(byt)
}
