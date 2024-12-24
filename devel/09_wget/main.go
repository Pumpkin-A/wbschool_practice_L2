package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

func main() {
	// uri := "https://sun9-7.userapi.com/impg/P5qwKdi3VEJvVXv6KvTrTs4ij9s9hHeyAMWr6w/d0lgFdt9UOY.jpg?size=806x1080&quality=96&sign=8d98653ac22de3ebf2b2f7864776df88&type=album"
	// uri := "http://example.com/download.mp4"
	// uri := "https://www.iana.org/domains/reserved"

	if len(os.Args) != 2 {
		fmt.Println("URI as argc should be provided")
		return
	}
	uri := os.Args[1]
	c := newDownloadCounter()
	c.download(uri, 1)
}

type counter struct {
	names map[string]struct{}
	m     *sync.Mutex
}

func newDownloadCounter() *counter {
	return &counter{names: make(map[string]struct{}), m: &sync.Mutex{}}
}

func (c *counter) isNew(uri string) bool {
	c.m.Lock()
	defer c.m.Unlock()
	if _, found := c.names[uri]; found {
		return false
	}
	c.names[uri] = struct{}{}
	return true
}

func (c *counter) download(uri string, maxLevel int) {
	if maxLevel < 0 {
		return
	}
	if _, e := url.ParseRequestURI(uri); e != nil {
		return
	}
	fmt.Println("dowload: ", uri, maxLevel)

	if !c.isNew(uri) {
		return
	}

	res, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
		return

	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	filename := parseFilename(uri, res.Header)
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	defer out.Close()

	out.Write(data)

	links := getSubLinks(data)
	var wg sync.WaitGroup
	for _, sublink := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			c.download(link, maxLevel-1)
		}(sublink)
	}
	wg.Wait()
}

func getSubLinks(data []byte) []string {
	re := regexp.MustCompile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	result := re.FindAll(data, -1)

	subUris := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		subUris = append(subUris, string(result[i]))
	}
	return subUris
}

func getFilename(url string, mediaType string) string {
	n := path.Base(url)

	if mediaType == "" {
		return "error-name"
	}
	name := cutAfter(cutAfter(n, "#"), "?")

	if path.Ext(name) == "" && mediaType != "" {
		return name + "." + mediaType
	}

	return name
}

func parseFilename(uri string, header http.Header) string {
	contentType := header.Get("Content-Type")
	mimeType, _, _ := mime.ParseMediaType(contentType)
	mediaType := cutBefore(mimeType, "/")

	filename := getFilename(uri, mediaType)

	return filename
}

func cutAfter(s, sep string) string {
	if strings.Contains(s, sep) {
		return strings.Split(s, sep)[0]
	}

	return s
}

func cutBefore(s, sep string) string {
	if strings.Contains(s, sep) {
		return strings.Split(s, sep)[1]
	}

	return s
}
