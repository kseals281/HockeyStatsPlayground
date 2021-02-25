package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func main() {
	resp := getPage("https://www.hockeydb.com/ihdb/stats/leagues/seasons/teams/0000332021.html")
	pprint(parseTable(resp))
}

func getPage(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		_ = fmt.Errorf("unable to retrieve page: %v", err)
	}

	return resp
}

func parseTable(resp *http.Response) [][]string {
	z := html.NewTokenizer(resp.Body)
	defer resp.Body.Close()
	var content [][]string
	var curr []string

	for {
		tt := z.Next()

		switch tt {

		case html.ErrorToken:
			return content

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "td" || t.Data == "th" {
				curr = append(curr, parseCell(z))
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "tr" {
				if len(curr) > 0 {
					content = append(content, curr)
					curr = []string{}
				}
			}

		}
	}
}

func parseCell(z *html.Tokenizer) string {
	inner := z.Next()
	if inner == html.TextToken {
		text := (string)(z.Text())
		t := strings.TrimSpace(text)
		return t
	}
	return parseCell(z)
}

func pprint(s [][]string) {
	for _, player := range s {
		fmt.Println(player)
	}
}
