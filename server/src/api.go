package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type htmlParseResponse struct {
	HtmlVersion            string `json:"html_version"`
	Title                  string `json:"title"`
	H1Count                int    `json:"h1_count"`
	H2Count                int    `json:"h2_count"`
	H3Count                int    `json:"h3_count"`
	H4Count                int    `json:"h4_count"`
	H5Count                int    `json:"h5_count"`
	H6Count                int    `json:"h6_count"`
	InternalLinksCount     int    `json:"internal_links_count"`
	ExternalLinksCount     int    `json:"external_links_count"`
	InaccessibleLinksCount int    `json:"inaccessible_links_count"`
	HasLoginForm           bool   `json:"has_login_form"`
}

func escapeUrl(url_string string) (string, error) {
	url_string = strings.Trim(url_string, " ")
	URL, err := url.Parse(url_string)
	if err != nil {
		return "", err
	}
	q := URL.Query()
	URL.RawQuery = q.Encode()
	return URL.String(), nil
}

func getData(c *gin.Context) {
	url_string := c.Query("url")
	URL, err := escapeUrl(url_string)
	if err != nil {
		fmt.Println("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := http.Get(URL)
	if err != nil {
		fmt.Println("error", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code err %d", res.StatusCode)
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		fmt.Println("error: ", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	data := parseHtml(doc, url_string)
	c.JSON(http.StatusOK, data)
}

func parseHtml(doc *html.Node, host string) (parsedData htmlParseResponse) {
	var urls []string
	var recParse func(*html.Node)
	var passwordCount int
	recParse = func(node *html.Node) {
		if node.Type == html.ElementNode {

			switch node.Data {
			case "html":
			case "h1":
				parsedData.H1Count++
			case "h2":
				parsedData.H2Count++
			case "h3":
				parsedData.H3Count++
			case "h4":
				parsedData.H4Count++
			case "h5":
				parsedData.H5Count++
			case "input":

				for _, a := range node.Attr {
					if a.Key == "type" && a.Val == "password" && isFormParent(node) {
						parsedData.HasLoginForm = true
					}
				}
			case "h6":
				parsedData.H6Count++
			case "a":
				for _, a := range node.Attr {
					if a.Key == "href" {
						//fmt.Println(a.Val)
						urls = append(urls, a.Val)
						break
					}
				}
			}

		} else if node.Type == html.TextNode {
			// reddit had multiple <title> tags so checked for title whos parent is head
			if node.Parent.Data == "title" && node.Parent.Type == html.ElementNode && node.Parent.Parent.Data == "head" {
				parsedData.Title = node.Data
			}
		} else if node.Type == html.DoctypeNode {
			for _, at := range node.Attr {
				switch at.Key {
				case "public":
					fallthrough
				case "system":
					parsedData.HtmlVersion += at.Val + " "
				}
			}
			parsedData.HtmlVersion = strings.TrimSpace(parsedData.HtmlVersion)
		}
		if node.Data == "input" {
			passwordCount++
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			recParse(c)
		}
	}
	recParse(doc)
	// default html version
	if parsedData.HtmlVersion == "" {
		parsedData.HtmlVersion = "HTML5"
	}
	parsedData.InternalLinksCount, parsedData.ExternalLinksCount,
		parsedData.InaccessibleLinksCount = getLinksData(urls, host)
	return
}

// recursive fn to check if <form> lies in the path to <html> from <input type="password">
func isFormParent(node *html.Node) bool {
	if node.Data == "html" {
		return false
	}
	if node.Data == "form" {
		return true
	}
	if node.Parent == nil {
		return false
	}
	return isFormParent(node.Parent)
}

func getLinksData(urls []string, host string) (internalLinksCount, externalLinksCount, inaccessibleLinksCount int) {
	host = strings.TrimSuffix(host, "/")
	hostUrl, err := url.Parse(host)
	validUrls := []string{}
	if err != nil {
		return
	}

	for _, url_string := range urls {
		url, err := url.Parse(url_string)
		if err != nil {
			fmt.Println("error parsing url", err)
			continue
		} else if len(url_string) == 0 {
			inaccessibleLinksCount++
			continue
		}
		//	fmt.Println(url_string, "HOST: ", url.Host, "URL_STRING_LEN: ", len(url_string))
		// subdomain urls are treated as external [https://stackoverflow.com/a/58289351]
		if url.Host == "" { // can start with / or letter
			internalLinksCount++
			if string(url_string[0]) == "/" {
				url_string = fmt.Sprintf("%s://%s%s", hostUrl.Scheme, hostUrl.Host, url_string)
			} else {
				url_string = fmt.Sprintf("%s://%s/%s", hostUrl.Scheme, hostUrl.Host, url_string)
			}
		} else if hostUrl.Host == url.Host {
			internalLinksCount++
		} else {
			externalLinksCount++
			if url.Scheme == "" { //example : //twitter.com/vanshajgirotra
				url_string = fmt.Sprintf("%s:%s", hostUrl.Scheme, url_string)
			}
		}
		url_string, err := escapeUrl(url_string)
		if err != nil {
			inaccessibleLinksCount++
			continue
		}
		validUrls = append(validUrls, url_string)
	}
	inaccessibleLinksCount = checkUrls(validUrls)
	return
}

func checkUrls(urls []string) (inaccessableLinkCount int) {
	c := make(chan bool)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go checkUrl(url, c, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	for isAccessable := range c {
		if !isAccessable {
			inaccessableLinkCount++
		}
	}
	return
}

func checkUrl(url string, c chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Head(url)
	// should use 3XX status too ?
	if err != nil || !(res.StatusCode >= 200 && res.StatusCode < 400) {
		c <- false
	} else {
		c <- true
	}
}
