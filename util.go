package main

import (
	"net/url"
	"path"
	"strings"
)

// URLJoinPath backport from go 1.19, code combined from:

// https://cs.opensource.google/go/go/+/refs/tags/go1.19beta1:src/net/url/url.go;l=1252
// https://cs.opensource.google/go/go/+/refs/tags/go1.19beta1:src/net/url/url.go;l=1252;bpv=1;bpt=1

// JoinPath returns a new URL with the provided path elements joined to
// any existing path and the resulting path cleaned of any ./ or ../ elements.
// Any sequences of multiple / characters will be reduced to a single /.
func URLJoinPath(base string, elem ...string) (result string, err error) {
	u, err := url.Parse(base)
	if err != nil {
		return
	}
	if len(elem) > 0 {
		elem = append([]string{u.EscapedPath()}, elem...)
		p := path.Join(elem...)
		// path.Join will remove any trailing slashes.
		// Preserve at least one.
		if strings.HasSuffix(elem[len(elem)-1], "/") && !strings.HasSuffix(p, "/") {
			p += "/"
		}
		u.Path = p
	}
	result = u.String()

	return
}

func parseURL(URL string) (u *url.URL, err error) {
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		// assume missing http[s] scheme by port
		host := strings.SplitN(URL, "/", 2)[0]
		hostPort := strings.SplitN(host, ":", 2)
		if len(hostPort) == 2 && hostPort[1] == "443" {
			URL = "https://" + URL
		} else {
			URL = "http://" + URL // default to http
		}
	}
	u, err = url.Parse(URL)
	// add default port, if missing
	hostPort := strings.SplitN(u.Host, ":", 2)
	if len(hostPort) != 2 {
		if u.Scheme == "http" {
			u.Host += ":80"
		} else if u.Scheme == "https" {
			u.Host += ":443"
		}
	}
	return
}

func parseURLWithRelativeHost(URL string, defaultHost string) (u *url.URL, err error) {
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		// assume missing http[s] scheme by port
		host := strings.SplitN(URL, "/", 2)[0]
		if len(host) == 0 {
			host = defaultHost
			if strings.HasPrefix(URL, "/") {
				URL = host + URL
			} else {
				URL = host + "/" + URL
			}
		}
		hostPort := strings.SplitN(host, ":", 2)
		if len(hostPort) == 2 && hostPort[1] == "443" {
			URL = "https://" + URL
		} else {
			URL = "http://" + URL // default to http
		}
	}
	u, err = url.Parse(URL)
	// add default port, if missing
	hostPort := strings.SplitN(u.Host, ":", 2)
	if len(hostPort) != 2 {
		if u.Scheme == "http" {
			// u.Host += ":80"
		} else if u.Scheme == "https" {
			// u.Host += ":443"
		}
	}
	return
}
