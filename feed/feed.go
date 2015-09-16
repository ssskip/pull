// Copyright 2015 The Pull Authors.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package feed


import (
	"encoding/xml"
	"html/template"
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"bytes"
	"io"
)


var (
	InValidFeedContentError = errors.New("invalid feed content")
	InValidOPMLContentError = errors.New("invalid opml content")
)

type Rss2 struct {
	XMLName     xml.Name    `xml:"rss"`
	Version     string        `xml:"version,attr"`
	Title       string        `xml:"channel>title"`
	Link        string        `xml:"channel>link"`
	Description string        `xml:"channel>description"`
	PubDate     string        `xml:"channel>pubDate"`
	ItemList    []Item        `xml:"channel>item"`
}


type Item struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Guid        string        `xml:"guid"`
	Description template.HTML    `xml:"description"`
	Content     template.HTML    `xml:"encoded"`
	PubDate     string        `xml:"pubDate"`
	Comments    string        `xml:"comments"`
}

type Link struct {
	Href string        `xml:"href,attr"`
}

type Atom1 struct {
	XMLName   xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Title     string        `xml:"title"`
	Subtitle  string        `xml:"subtitle"`
	Id        string        `xml:"id"`
	Updated   string        `xml:"updated"`
	Rights    string        `xml:"rights"`
	Link      Link        `xml:"link"`
	Author    Author        `xml:"author"`
	EntryList []Entry        `xml:"entry"`
}

type Author struct {
	Name  string        `xml:"name"`
	Email string        `xml:"email"`
}

type Entry struct {
	Title   string        `xml:"title"`
	Summary string        `xml:"summary"`
	Content string        `xml:"content"`
	Id      string        `xml:"id"`
	Updated string        `xml:"updated"`
	Link    Link        `xml:"link"`
	Author  Author        `xml:"author"`
}



// OPML is the root node of an OPML document. It only has a single required
// attribute: the version.
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head holds some meta information about the document.
type Head struct {
	Title           string `xml:"title"`
	DateCreated     string `xml:"dateCreated,omitempty"`
	DateModified    string `xml:"dateModified,omitempty"`
	OwnerName       string `xml:"ownerName,omitempty"`
	OwnerEmail      string `xml:"ownerEmail,omitempty"`
	OwnerID         string `xml:"ownerId,omitempty"`
	Docs            string `xml:"docs,omitempty"`
	ExpansionState  string `xml:"expansionState,omitempty"`
	VertScrollState string `xml:"vertScrollState,omitempty"`
	WindowTop       string `xml:"windowTop,omitempty"`
	WindowBottom    string `xml:"windowBottom,omitempty"`
	WindowLeft      string `xml:"windowLeft,omitempty"`
	WindowRight     string `xml:"windowRight,omitempty"`
}

// Body is the parent structure of all outlines.
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline holds all information about an outline.
type Outline struct {
	Outlines     []Outline `xml:"outline"`
	Text         string    `xml:"text,attr"`
	Type         string    `xml:"type,attr,omitempty"`
	IsComment    string    `xml:"isComment,attr,omitempty"`
	IsBreakpoint string    `xml:"isBreakpoint,attr,omitempty"`
	Created      string    `xml:"created,attr,omitempty"`
	Category     string    `xml:"category,attr,omitempty"`
	XMLURL       string    `xml:"xmlUrl,attr,omitempty"`
	HTMLURL      string    `xml:"htmlUrl,attr,omitempty"`
	URL          string    `xml:"url,attr,omitempty"`
	Language     string    `xml:"language,attr,omitempty"`
	Title        string    `xml:"title,attr,omitempty"`
	Version      string    `xml:"version,attr,omitempty"`
	Description  string    `xml:"description,attr,omitempty"`
}


func atom1ToRss2(a Atom1) Rss2 {
	r := Rss2{
		Title: a.Title,
		Link: a.Link.Href,
		Description: a.Subtitle,
		PubDate: a.Updated,
	}
	r.ItemList = make([]Item, len(a.EntryList))
	for i, entry := range a.EntryList {
		r.ItemList[i].Title = entry.Title
		r.ItemList[i].Link = entry.Link.Href
		if entry.Content == "" {
			r.ItemList[i].Description = template.HTML(entry.Summary)
		} else {
			r.ItemList[i].Description = template.HTML(entry.Content)
		}
	}
	return r
}

func parseXML(content []byte, v interface{}) error {
	d := xml.NewDecoder(bytes.NewReader(content))
	d.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		//converts GBK to UTF-8.
		if s == "GBK" {
			return transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder()), nil
		}
		return charset.NewReader(r, s)
	}
	err := d.Decode(v)
	return err
}

func parseAtom1ToRss2(content []byte) (Rss2, error) {
	a := Atom1{}
	err := parseXML(content, &a)
	if err != nil {
		return Rss2{}, err
	}
	return atom1ToRss2(a), nil
}


func ParseFeedContent(content []byte) (Rss2, error) {
	r := Rss2{}
	err := parseXML(content, &r)
	if err != nil {
		return parseAtom1ToRss2(content)
	}

	if r.Version == "2.0" {
		for i, _ := range r.ItemList {
			if r.ItemList[i].Content != "" {
				r.ItemList[i].Description = r.ItemList[i].Content
			}
		}
		return r, nil
	}
	return r, InValidFeedContentError
}
func ParseOPMLContent(content []byte) (OPML, error) {
	o := OPML{}
	err := parseXML(content, &o)
	if err != nil {
		return o, InValidOPMLContentError
	}
	return o, nil
}
