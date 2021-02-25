package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func main() {
	resp := getPage("https://www.hockeydb.com/ihdb/stats/leagues/seasons/teams/0000332021.html")
	parseTable(resp)
}

func getPage(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		_ = fmt.Errorf("unable to retrieve page: %v", err)
	}

	return resp
}

func parseTable(resp *http.Response) {
	z := html.NewTokenizer(resp.Body)
	defer resp.Body.Close()
	var content [][]string
	var curr []string

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			pprint(content)
			return

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "td" {
				inner := z.Next()
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := strings.TrimSpace(text)
					curr = append(curr, t)
				}
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

func pprint(s [][]string) {
	for _, player := range s {
		fmt.Println(player)
	}
}
