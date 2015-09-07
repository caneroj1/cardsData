package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// PostToDB sends a post request via curl to my local
// cardsAPI to POST a new card
func PostToDB() {
	if len(os.Args) < 3 {
		fmt.Println("Incorrect number of arguments supplied. The first argument after post must be the required body of the card.")
		os.Exit(1)
	}

	postURL := "http://localhost:9000/api/cards"
	fmt.Printf("POSTING to %s\n", postURL)
	for _, arg := range os.Args[2:] {
		fmt.Printf("\t- %s\n", arg)
	}

	var (
		cardType, blanks string
	)

	body := os.Args[2]
	if len(os.Args) > 3 {
		cardType = os.Args[3]
	}
	if len(os.Args) > 4 {
		blanks = os.Args[4]
	}

	var client http.Client
	data := make(url.Values)
	data["cardBody"] = []string{body}
	data["cardType"] = []string{cardType}
	data["cardBlanks"] = []string{blanks}
	resp, err := client.PostForm(postURL, data)

	if err != nil {
		fmt.Println(err)
	}

	code := resp.StatusCode
	fmt.Printf("Status: %d\n", code)
	fmt.Println(readBody(resp.Body))

	// if code == 200 {
	// 	fmt.Println(readBody(resp.Body))
	// } else {
	// 	fmt.Println(readErrors(resp.Cookies()))
	// }
}

func readBody(body io.ReadCloser) string {
	bytes := make([]byte, 100)
	var returnStr string
	for {
		n, err := body.Read(bytes)
		if err != nil {
			if err == io.EOF {
				if n > 0 {
					returnStr += string(bytes[:n])
				}
				return returnStr
			}
		}
		returnStr += string(bytes[:n])
	}
}

func readErrors(cookies []*http.Cookie) string {
	fmt.Println("Errors:")
	for _, cookie := range cookies {
		if cookie.Name == "REVEL_ERRORS" {
			str, err := url.QueryUnescape(cookie.String())
			if err != nil {
				fmt.Println(err)
			} else {
				return str
			}
		}
	}

	return ""
}
