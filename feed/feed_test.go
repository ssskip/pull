// Copyright 2015 The Pull Authors.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package feed_test

import (
	"io/ioutil"
	"testing"
	"pull/feed"
)

var (
	rss2, _ = ioutil.ReadFile("testdata/rss2.xml")
	atom1, _ = ioutil.ReadFile("testdata/atom1.xml")
	opml, _ = ioutil.ReadFile("testdata/opml.xml")
	feedburner, _ = ioutil.ReadFile("testdata/feedburner.xml")
)

func TestParseFeed(t *testing.T) {
	rss, err := feed.ParseFeedContent(rss2)
	if err != nil {
		t.Error(err)
	}
	t.Log(rss.Title, rss.Link, rss.ItemList[0].PubDate)
	rss, err = feed.ParseFeedContent(atom1)
	if err != nil {
		t.Error(err)
	}
	t.Log(rss.Title, rss.Link)
}

func TestFeedBurner(t *testing.T) {
	rss, err := feed.ParseFeedContent(feedburner)
	if err != nil {
		t.Error(err)
	}
	t.Log(rss.Title, rss.Link, rss.ItemList[0].PubDate)
}
func TestParseOPML(t *testing.T) {
	opml, err := feed.ParseOPMLContent(opml)
	if err != nil {
		t.Error(err)
	}
	t.Log(opml.Head.Title)
}