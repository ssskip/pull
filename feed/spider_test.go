// Copyright 2015 The Pull Authors.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package feed_test

import (
	"testing"
	gourl "net/url"
	"pull/feed"
)

func TestFeedSpider(t *testing.T) {
	var result *feed.Response
	url, _ := gourl.Parse("https://www.percona.com/blog/feed/atom/")
	result = <-(&feed.Spider{
		Url:url,
		Timeout:30,
	}).Run()
	if result.Err != nil {
		t.Error(result.Err)
	}
	t.Log(result.Duration, result.Md5)
}