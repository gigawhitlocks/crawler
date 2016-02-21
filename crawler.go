package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
)

var url string

func init() {
	flag.Parse()
	if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	} else if len(flag.Args()) == 0 {
		url = "http://news.ycombinator.com"
		return
	}
	url = flag.Args()[0]
}

func ParseLinks(bodyHtml io.Reader) []string {
	var output []string
	t := html.NewTokenizer(bodyHtml)
	for {
		tt := t.Next()
		if tt == html.ErrorToken {
			return output
		}
		if tn, _ := t.TagName(); string(tn) == "a" &&
			tt == html.StartTagToken {

			key, val, _ := t.TagAttr()
			sval := string(val)

			if string(key) == "href" &&
				!(strings.HasPrefix(sval, "http") ||
					strings.HasPrefix(sval, "//")) {
				sval = url + "/" + sval
			}
			output = append(output, sval)
		}
	}
}

func main() {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	links := ParseLinks(response.Body)
	for _, link := range links {
		fmt.Println(link)
	}
}
