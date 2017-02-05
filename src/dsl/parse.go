package dsl

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type DSLScheme string

const (
	DSLPosixScheme      DSLScheme = "content"
	DSLFtpScheme        DSLScheme = "ftp"
	DSLPosixDefaultPort           = 5472
)

type DSLUrl struct {
	Host   string
	Scheme DSLScheme
	Path   string
	Port   int
	Url    *url.URL
}

func (d *DSLUrl) PosixAddress() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}

func (d *DSLUrl) ToString() string {
	return d.Url.String()
}

func parseUrl(in string) (*DSLUrl, error) {

	u, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	result := &DSLUrl{Host: u.Host, Path: strings.TrimLeft(u.Path, "/"), Url: u}
	// validate scheme
	switch string(u.Scheme) {

	// "posix"
	case string(DSLPosixScheme):
		result.Scheme = DSLPosixScheme
		if host, port, err := net.SplitHostPort(u.Host); err == nil {
			result.Host = host
			if parsedPort, err := strconv.ParseInt(port, 10, 64); err == nil {
				result.Port = int(parsedPort)
			}
		} else {
			result.Port = DSLPosixDefaultPort
		}

	// ftp
	case string(DSLFtpScheme):
		result.Scheme = DSLFtpScheme

	default:
		return nil, fmt.Errorf("parse url %s unknown scheme: %s", u.String(), u.Scheme)

	}

	return result, nil
}
