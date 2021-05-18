package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func TestCheckUrls(t *testing.T) {
	cases := []struct {
		urls []string
		want int
	}{
		{[]string{"https://example.com"}, 0},
		{[]string{"http://example.com"}, 0},
		{[]string{"example.com"}, 1}, // no protocol test
		{[]string{"https://thisurldoesnotexist.com/"}, 1},
		{[]string{"https://abcdlaksjdslkf.com/"}, 1},
		{[]string{"http://abcdlaksjdslkf.com/"}, 1},
		{[]string{"abcdjasdasd"}, 1},
		{[]string{"mailto:abcd@gmail.com"}, 1},
		{[]string{"tel:19287129847"}, 1},
		{[]string{"https://myaccount.google.com/"}, 0},
		{[]string{}, 0},
	}

	for _, test := range cases {
		fmt.Println("testing", test.urls)
		got := checkUrls(test.urls)
		if got != test.want {
			t.Errorf("checkUrls(%q), want %v got %v", test.urls, test.want, got)
		}
	}
}

func TestGetLinksData(t *testing.T) {
	type linkdata struct {
		InternalLinksCount     int
		ExternalLinksCount     int
		InaccessibleLinksCount int
	}
	cases := []struct {
		urls         []string
		requestedUrl string
		want         linkdata
	}{
		{[]string{"https://example.com/?hello", "/", "https://google.com"}, "https://example.com", linkdata{2, 1, 0}},
		{[]string{"https://youtube.com", "https://reddit.com/login", "https://google.com"}, "https://example.com", linkdata{0, 3, 0}},
		{[]string{"/register/", "https://redditinc.com", "/username", "reddit.com/signup", "https://example.com", "example.com"},
			"https://reddit.com/login", linkdata{4, 2, 2}},
		{[]string{"thisdoesnotexist", "https://stackoverflow.com/users/", "https://linkedin.com", "/12333", "resume", "/resume"},
			"https://www.vanshajgirotra.com", linkdata{4, 2, 2}},
		{[]string{}, "https://www.vanshajgirotra.com", linkdata{0, 0, 0}},
		{[]string{"thisdoesnotexist", "https://stackoverflow.com/users/", "https://linkedin.com", "/12333", "resume", "/resume"},
			"https://www.vanshajgirotra.com", linkdata{4, 2, 2}},
		{[]string{"https://myaccount.google.com/", "https://mail.google.com/", "https://google.org"}, "https://google.com", linkdata{0, 3, 0}},
		{[]string{"//twitter.com/vanshajgirotra", "https://github.com/vanshajg", "/hello.com"}, "https://www.vanshajgirotra.com", linkdata{1, 2, 1}},
	}
	for _, test := range cases {
		fmt.Println("testing", test.urls)
		internalLinksCount, externalLinksCount, inaccessibleLinksCount := getLinksData(test.urls, test.requestedUrl)
		got := linkdata{internalLinksCount, externalLinksCount, inaccessibleLinksCount}

		if got != test.want {
			t.Errorf("checkUrls(%q), want %v got %v", test.urls, test.want, got)
		}
	}
}

func TestParseHtml(t *testing.T) {

	cases := []struct {
		path string
		host string
		want htmlParseResponse
	}{
		{
			"../test_html/only_password.html",
			"https://example.com",
			htmlParseResponse{
				"HTML5",
				"",
				0, 0, 0, 0, 0, 0,
				0, 0, 0,
				false,
			},
		},
		{
			"../test_html/basic_form.html",
			"https://example.com",
			htmlParseResponse{
				"HTML5",
				"",
				0, 0, 0, 0, 0, 0,
				0, 0, 0,
				true,
			},
		},
		{
			"../test_html/only_password.html",
			"https://example.com",
			htmlParseResponse{
				"HTML5",
				"",
				0, 0, 0, 0, 0, 0,
				0, 0, 0,
				false,
			},
		},
		{
			"../test_html/link_test.html",
			"https://example.com",
			htmlParseResponse{
				"HTML5",
				"",
				0, 0, 0, 0, 0, 0,
				2, 1, 1,
				true,
			},
		},
		{
			"../test_html/heading_test.html",
			"https://example.com",
			htmlParseResponse{
				"HTML5",
				"",
				9, 2, 2, 1, 0, 0,
				0, 0, 0,
				false,
			},
		},
		{
			"../test_html/mixed_1.html",
			"https://example.com",
			htmlParseResponse{
				"HTML5",
				"Hello World",
				1, 0, 0, 0, 0, 0,
				1, 2, 1,
				false,
			},
		},
		{
			"../test_html/mixed_2.html",
			"https://example.com",
			htmlParseResponse{
				"-//W3C//DTD HTML 4.01//EN http://www.w3.org/TR/html4/strict.dtd",
				"html 4 page",
				1, 1, 1, 1, 0, 0,
				0, 0, 0,
				true,
			},
		},
	}

	for _, test := range cases {
		func() {
			fmt.Println("testing", test.path)
			data, err := os.Open(test.path)

			if err != nil {
				fmt.Println("error opening file", err)
				return
			}
			defer data.Close()

			html, err := html.Parse(data)
			if err != nil {
				fmt.Println("error parsing html", err)
				return
			}
			got := parseHtml(html, test.host)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("check html %v got: %v ; want: %v", test.path, got, test.want)
			}

		}()
	}

}
