package crawler

import (
	"testing"
)

var (
	fullProcessTests = []struct {
		input  string
		source string
		output string
	}{
		{
			"../a.html",
			"http://velikodny.com",
			"http://velikodny.com/a.html",
		},
		{
			"../a.html",
			"http://velikodny.com/a",
			"http://velikodny.com/a.html",
		},
		{
			"../a.html",
			"http://velikodny.com/a/",
			"http://velikodny.com/a.html",
		},
		{
			"../a.html",
			"http://velikodny.com/aaa/bbb",
			"http://velikodny.com/a.html",
		},
		{
			"../a.html",
			"http://velikodny.com/aaa/bbb/",
			"http://velikodny.com/aaa/a.html",
		},
		{
			"./a.html",
			"http://velikodny.com/aaa/bbb",
			"http://velikodny.com/aaa/a.html",
		},
		{
			"./a.html",
			"http://velikodny.com/aaa/bbb/",
			"http://velikodny.com/aaa/bbb/a.html",
		},
		{
			"/a.html",
			"http://velikodny.com/",
			"http://velikodny.com/a.html",
		},
		{
			"#",
			"http://velikodny.com/",
			"http://velikodny.com/",
		},
		{
			"#aaa",
			"http://velikodny.com/",
			"http://velikodny.com/",
		},
		{
			"?",
			"http://velikodny.com/",
			"http://velikodny.com/",
		},
		{
			"?aa",
			"http://velikodny.com/",
			"http://velikodny.com/?aa=",
		},
		{
			"?aa",
			"http://velikodny.com/",
			"http://velikodny.com/?aa=",
		},
		{
			"//www.velikodny.com/",
			"https://www.velikodny.com/news/",
			"https://www.velikodny.com/",
		},
	}
)

func TestLinkProcessor2(t *testing.T) {
	processor := &linkProcessor{domains: []string{"velikodny.com", "www.velikodny.com"}}

	l := &Link{Ref: "http://#", Source: "http://velikodny.com"}
	processor.Process(l)

	if !l.Malformed {
		t.Fatalf("should be malformed %+v", l)
	}

	l = &Link{Ref: "http://", Source: "http://velikodny.com"}
	processor.Process(l)

	if !l.Malformed {
		t.Fatalf("should be malformed %+v", l)
	}
}

func TestFullLinkProcess(t *testing.T) {
	processor := &linkProcessor{domains: []string{"velikodny.com", "www.velikodny.com"}}

	for tn, test := range fullProcessTests {
		l := &Link{Ref: test.input, Source: test.source}
		processor.Process(l)

		if l.Ref != test.output {
			t.Fatalf("%v: expected '%s', but got '%v'", tn, test.output, l.Url.String())
		}
	}
}
