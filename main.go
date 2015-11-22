package search_result_crawler

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

type Result struct {
	Title       string
	URL         string
	Description string
}

type IntermediatePair struct {
	Keyword string
	Page    int
	Index   int
	Result
}

type Parser interface {
	GetSearchResultPageURL(string, int) string                   // (keyword string, page int) -> (string)
	Parse(*goquery.Document, string, int, chan IntermediatePair) // (doc *goquery.Document, keyword string, page int, channel chan intermediatePair)
}

type Crawler struct {
	Input map[string][]int
	Parser
}

// input: map[string][]int
//           (keyword)(a slice of the page number you want to crawl)
func (c *Crawler) Start() map[string][][]Result {
	channel := make(chan IntermediatePair)

	pageTotal := 0
	for _, pagesNeedToBeCrawled := range c.Input {
		pageTotal += len(pagesNeedToBeCrawled)
	}

	wg.Add(pageTotal)
	go c.deliverSearchResultParsingTasks(channel)
	var result map[string][][]Result
	result = reduceResult(channel)

	return result
}

func (c *Crawler) deliverSearchResultParsingTasks(channel chan IntermediatePair) {
	for keyword, pagesNeedToBeCrawled := range c.Input {
		for _, page := range pagesNeedToBeCrawled {
			reader := fetchDocument(c.Parser.GetSearchResultPageURL(keyword, page))
			doc, err := goquery.NewDocumentFromReader(reader)

			if err != nil {
				break
			}

			go func(doc *goquery.Document, keyword string, page int, channel chan IntermediatePair) {
				c.Parser.Parse(doc, keyword, page, channel)
				wg.Done()
			}(doc, keyword, page, channel)
		}

	}
	wg.Wait()
	close(channel)
}

func reduceResult(channel chan IntermediatePair) map[string][][]Result {
	var result map[string][][]Result = make(map[string][][]Result)
	//  (keyword map)(page slice)(value string)
	// keyword: {page1: {}, page2: {}}

	for pair := range channel {
		for pair.Page+1 > len(result[pair.Keyword]) {
			result[pair.Keyword] = append(result[pair.Keyword], make([]Result, 0))
		}
		for pair.Index+1 > len(result[pair.Keyword][pair.Page]) {
			result[pair.Keyword][pair.Page] = append(result[pair.Keyword][pair.Page], Result{})
		}
		result[pair.Keyword][pair.Page][pair.Index] = pair.Result
	}
	return result
}

func fetchDocument(url string) io.Reader {
	buffer := bytes.NewBufferString("")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.86 Safari/537.36")
	res, _ := http.DefaultClient.Do(req)
	res.Write(buffer)
	defer res.Body.Close()
	return buffer
}
