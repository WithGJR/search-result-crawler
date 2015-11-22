# A Framework For Writing Your Own Search Result Crawler in Golang

## Installation

``` bash
go get github.com/WithGJR/search-result-crawler
```

## How to Use

First, you need to initialize a `Crawler` struct. `Crawler` struct has two fields, one is `Input`, and another one is `Parser`. 

`Parser` is an interface, the structure of this interface is:

``` golang
type Parser interface {
	GetSearchResultPageURL(string, int) string                   // (keyword string, page int) -> (string)
	Parse(*goquery.Document, string, int, chan IntermediatePair) // (doc *goquery.Document, keyword string, page int, channel chan intermediatePair)
}
```

So, you need to implement this interface by yourself.

After initialization, you just need to call the `Crawler`'s method `Start()`, this method will start the crawling tasks and return the parsing `result`. 

The type of `result` is `map[string][][]Result`. The key of this map is the keyword you specified, and the value of this map is a slice of search result pages of a specified keyword, and the elements of this slice is another slice of `Result`.

The structure of `Result` is:

``` golang
type Result struct {
  Title       string
  URL         string
  Description string
}
```

Following is a simple example showing how to use this project:

``` golang
  c := crawler.Crawler{
    Input: map[string][]int{
      "golang tutorial": []int{1, 2, 3},
      "seo":             []int{1, 2, 3, 5},
    },
    Parser: &parser.GoogleSearchParser{},
  }
  r := c.Start()
  // If you call r["seo"][0], you will get a slice of the first page of the search result of the keyword 'seo'.
```
## Parser Example

- [Google Taiwan Search Parser](https://github.com/WithGJR/google-search-parser)

## Chinese Tutorial

I wrote a Chinese tutorial and posted it on my blog: [Golang - 如何寫爬蟲抓搜尋引擎的搜尋結果？](http://blog.cgmlife.net/posts/2015/11/22/how-to-write-search-result-crawler-in-golang).
