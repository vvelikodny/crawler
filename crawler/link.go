package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"
)

var (
	ErrRefEmpty = errors.New("reference should not be empty")
)

type Link struct {
	Source    string `bson:"Source"`
	RawRef    string `bson:"RawRef"`
	Ref       string `bson:"Ref"`
	Url       *url.URL
	Malformed bool   `bson:"Malformed"`
	SelfLink  bool   `bson:"SelfLink"`
	Error     string `bson:"Error"`
}

func NewLink(ref string) (*Link, error) {
	if len(ref) == 0 {
		return nil, ErrRefEmpty
	}

	return &Link{RawRef: ref, Ref: ref}, nil
}

func NewHrefLink(source *Link, href string) *Link {
	return &Link{Source: source.Ref, RawRef: href, Ref: href}
}

func (l *Link) Md5() string {
	md5v := md5.Sum([]byte(l.Ref))
	return hex.EncodeToString(md5v[:])
}

func (l *Link) IsRejected() bool {
	return l.Malformed
}

func (l *Link) SetMalformed(error string) {
	l.Malformed = true
	l.Error = error
}
