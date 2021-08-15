package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gaorx/stardust3/sdresty"
)

func main() {
	c := sdresty.New(sdresty.Options{})
	resp, err := c.R().Get("http://www.useragentstring.com/pages/useragentstring.php?name=All")
	if err != nil {
		panic(err)
	}

	html := resp.String()
	if err != nil {
		panic(err)
	}

	patt := regexp.MustCompile(`<a href='/index.php\?id=\d+'>([^<]+)</a>`)
	l := patt.FindAllStringSubmatch(html, -1)
	var lines []string
	for _, ss := range l {
		ua := ss[1]
		lines = append(lines, fmt.Sprintf(`		%s,`, strconv.Quote(ua)))
	}

	t := `
package sdhttpua

var (
	rawUserAgents = []string{
%s
	}
)
`
	goFile := fmt.Sprintf(t, strings.Join(lines, "\n"))
	fmt.Println(goFile)
}
