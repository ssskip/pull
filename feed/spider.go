// Copyright 2015 The Pull Authors.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package feed

import (
	"net/url"
	"time"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"crypto/md5"
	"fmt"
)

const AGENT_NAME = "PULL RSS Reader " + VERSION

type Response struct {
	Err        error
	StatusCode int
	Rss        *Rss2
	Md5        string
	Duration   time.Duration
}


type Spider struct {
	Url     *url.URL
	// Timeout in seconds.
	Timeout int
}

func (s *Spider) Run() (<-chan *Response) {
	out := make(chan *Response)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
		DisableKeepAlives:  false,
		// TODO Append dial timeout.
		TLSHandshakeTimeout: time.Duration(s.Timeout) * time.Second,
	}}
	req, _ := http.NewRequest("GET", s.Url.String(), nil)
	req.Header.Set("User-Agent", AGENT_NAME)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en")

	go func() {
		start := time.Now()
		code := 0
		resp, err := client.Do(req)
		if err != nil {
			out <- (&Response{
				Err:err,
			})
			return
		}
		code = resp.StatusCode
		content, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			out <- (&Response{
				Err:err,
				StatusCode:code,
			})
			return
		}
		rss, err := ParseFeedContent(content)
		out <- (&Response{
			Err:err,
			StatusCode:code,
			Rss:&rss,
			Md5:fmt.Sprintf("%x", md5.Sum(content)),
			Duration:time.Now().Sub(start),
		})
		close(out)
	}()

	return out
}