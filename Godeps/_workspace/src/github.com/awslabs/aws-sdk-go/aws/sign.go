package aws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

// Context encapsulates the context of a client's connection to an AWS service.
type Context struct {
	Service     string
	Region      string
	Credentials CredentialsProvider
}

func (c *Context) sign(r *http.Request) error {
	date := r.Header.Get("Date")
	t := currentTime().UTC()
	if date != "" {
		var err error
		t, err = time.Parse(http.TimeFormat, date)
		if err != nil {
			return err
		}
	}
	r.Header.Set("x-amz-date", t.Format(iso8601BasicFormat))

	chash, err := c.hashContent(r)
	if err != nil {
		return err
	}
	r.Header.Set("x-amz-content-sha256", chash)

	creds, err := c.Credentials.Credentials()
	if err != nil {
		return err
	}

	if s := creds.SessionToken; s != "" {
		r.Header.Set("X-Amz-Security-Token", s)
	}

	k := c.signature(creds.SecretAccessKey, t)
	h := hmac.New(sha256.New, k)
	c.writeStringToSign(h, t, r, chash)

	auth := bytes.NewBufferString("AWS4-HMAC-SHA256 ")
	_, _ = auth.WriteString("Credential=" + creds.AccessKeyID + "/" + c.creds(t))
	_, _ = auth.WriteString(", ")
	_, _ = auth.WriteString("SignedHeaders=")
	c.writeHeaderList(auth, r)
	_, _ = auth.WriteString(", ")
	_, _ = auth.WriteString(fmt.Sprintf("Signature=%x", h.Sum(nil)))

	r.Header.Set("Authorization", auth.String())
	return nil
}

func (c *Context) writeStringToSign(w io.Writer, t time.Time, r *http.Request, chash string) {
	_, _ = io.WriteString(w, "AWS4-HMAC-SHA256")
	_, _ = io.WriteString(w, lf)
	_, _ = io.WriteString(w, t.Format(iso8601BasicFormat))
	_, _ = io.WriteString(w, lf)

	_, _ = io.WriteString(w, c.creds(t))
	_, _ = io.WriteString(w, lf)

	h := sha256.New()
	c.writeRequest(h, r, chash)
	fmt.Fprintf(w, "%x", h.Sum(nil))
}

func (c *Context) writeRequest(w io.Writer, r *http.Request, chash string) {
	r.Header.Set("host", r.Host)

	_, _ = io.WriteString(w, r.Method)
	_, _ = io.WriteString(w, lf)
	c.writeURI(w, r)
	_, _ = io.WriteString(w, lf)
	c.writeQuery(w, r)
	_, _ = io.WriteString(w, lf)
	c.writeHeader(w, r)
	_, _ = io.WriteString(w, lf)
	_, _ = io.WriteString(w, lf)
	c.writeHeaderList(w, r)
	_, _ = io.WriteString(w, lf)
	fmt.Fprint(w, chash)
}

func (c *Context) hashContent(r *http.Request) (string, error) {
	var b []byte
	// If the payload is empty, use the empty string as the input to the SHA256 function
	// http://docs.amazonwebservices.com/general/latest/gr/sigv4-create-canonical-request.html
	if r.Body == nil {
		b = []byte("")
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		b = body
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}

	h := sha256.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (c *Context) writeURI(w io.Writer, r *http.Request) {
	p := r.URL.RequestURI()
	if r.URL.RawQuery != "" {
		p = p[:len(p)-len(r.URL.RawQuery)-1]
	}
	slash := strings.HasSuffix(p, "/")
	p = path.Clean(p)
	if p != "/" && slash {
		p += "/"
	}
	_, _ = io.WriteString(w, p)
}

func (c *Context) writeQuery(w io.Writer, r *http.Request) {
	var a []string
	for k, vs := range r.URL.Query() {
		k = url.QueryEscape(k)
		for _, v := range vs {
			if v == "" {
				a = append(a, k+"=")
			} else {
				v = url.QueryEscape(v)
				a = append(a, k+"="+v)
			}
		}
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			_, _ = io.WriteString(w, "&")
		}
		_, _ = io.WriteString(w, s)
	}
}

func (c *Context) writeHeader(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k, v := range r.Header {
		sort.Strings(v)
		a[i] = strings.ToLower(k) + ":" + strings.Join(v, ",")
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			_, _ = io.WriteString(w, lf)
		}
		_, _ = io.WriteString(w, s)
	}
}

func (c *Context) writeHeaderList(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k := range r.Header {
		a[i] = strings.ToLower(k)
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			_, _ = io.WriteString(w, ";")
		}
		_, _ = io.WriteString(w, s)
	}
}

func (c *Context) creds(t time.Time) string {
	return t.Format(iso8601BasicFormatShort) + "/" + c.Region + "/" + c.Service + "/aws4_request"
}

func (c *Context) signature(secretAccessKey string, t time.Time) []byte {
	h := ghmac(
		[]byte("AWS4"+secretAccessKey),
		[]byte(t.Format(iso8601BasicFormatShort)),
	)
	h = ghmac(h, []byte(c.Region))
	h = ghmac(h, []byte(c.Service))
	h = ghmac(h, []byte("aws4_request"))
	return h
}

func ghmac(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write(data)
	return h.Sum(nil)
}

var currentTime = time.Now

const (
	lf                      = "\n"
	iso8601BasicFormat      = "20060102T150405Z"
	iso8601BasicFormatShort = "20060102"
)
