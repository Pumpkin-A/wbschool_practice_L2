package main

import (
	"fmt"
	"io"
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
	// uri := "https://sun9-74.userapi.com/impg/EJoWbiBxQ4MmPXCcvtTGaAutkTnl3UGkO8_7gg/uouvsFrnkNE.jpg?size=1266x1583&quality=95&sign=5440a5d4cb91e3925784f2b1e389987f&type=album"
	// uri := "http://example.com/download.mp4"
	// uri := "https://www.iana.org/domains/reserved"

	if len(os.Args) != 2 {
		fmt.Println("URI as argc should be provided")
		return
	}
	uri := os.Args[1]
	c := newDownloadCounter()
	if err := c.download(uri, 1); err != nil {
		fmt.Println("err found")
	}
}

type counter struct {
	// Т.к. Names - разделяемый ресурс между горутинами, то синхронизируем доступ через мьютекс
	names map[string]struct{}
	m     *sync.Mutex
}

func newDownloadCounter() *counter {
	return &counter{
		names: make(map[string]struct{}),
		m:     &sync.Mutex{},
	}
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

// Основной метод работы
// - загружает данные
// - ищет ссылки в скачанных данных
// - скачивает данные по найденным ссылкам в горутнах
func (c *counter) download(uri string, maxLevel int) error {
	if maxLevel < 0 {
		return fmt.Errorf("maxlevel error")
	}
	if _, err := url.ParseRequestURI(uri); err != nil {
		return err
	}
	fmt.Println("dowload: ", uri, maxLevel)

	if !c.isNew(uri) {
		return nil
	}

	res, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
		return err

	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	filename := parseFilename(uri, res.Header)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := out.Write(data); err != nil {
		return err
	}

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

	return nil
}

func getSubLinks(data []byte) []string {
	// Вроде бы работает такой регекс для поиска ссылок
	re := regexp.MustCompile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	result := re.FindAll(data, -1)

	subUris := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		subUris = append(subUris, string(result[i]))
	}
	return subUris
}

// Создаем имя файла из ссылки медиа-данных о формате файла, чтобы его имя подходило и был правильный формат файла
func parseFilename(uri string, header http.Header) string {
	contentType := header.Get("Content-Type")
	mimeType, _, _ := mime.ParseMediaType(contentType)
	mediaType := cutBefore(mimeType, "/")

	filename := getFilename(uri, mediaType)

	return filename
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
