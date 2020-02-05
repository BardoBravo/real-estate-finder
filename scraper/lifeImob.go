package cloudrun

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type LifeImob struct{}

var LocalExporter Exporter

type Ext struct {
	*gocrawl.DefaultExtender
}

func (platform LifeImob) NewCollector(config Config) *colly.Collector {
	log.Print("Starting LifeImob")
	options := append(
		config.collectorOptions,
		colly.AllowedDomains("www.lifeimob.com.br"))
	return colly.NewCollector(options...)
}

func (platform LifeImob) crawl(config Config, exporter Exporter) *colly.Collector {

	LocalExporter = exporter

	colly := platform.NewCollector(config)

	ext := &Ext{&gocrawl.DefaultExtender{}}

	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	opts.MaxVisits = 999999999

	log.Print("starting crawl...")
	c := gocrawl.NewCrawlerWithOptions(opts)
	if err := c.Run("https://www.lifeimob.com.br/index.asp?id_pagina=8&OrdernarPor=Preco%20DESC&ValorDe=0&ValorAte=500000000&cidade=3326"); err != nil {
		log.Print(err)
	}
	return colly
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	spaceString := ""
	priceString := ""
	roomString := ""
	title := ""
	locationString := ""

	header := doc.Find(".texts-center")
	if header.Length() > 0 {
		header.Each(func(i int, n *goquery.Selection) {
			textHeader := n.Find("b")
			title = textHeader.Text()
			locationString = n.Find("h1").Text()
		})
	}

	data := doc.Find(".icons-dd")
	if data.Length() > 0 {
		data.Each(func(i int, n *goquery.Selection) {
			n.Find("li").Each(func(cont int, found *goquery.Selection) {
				if strings.Contains(found.Text(), "Área Total") {
					spaceString = found.Text()
				} else if strings.Contains(found.Text(), "Área Terreno") {
					spaceString = found.Text()
				} else if strings.Contains(found.Text(), "Área privativa") {
					spaceString = found.Text()
				} else if strings.Contains(found.Text(), "Dormitório(s)") {
					roomString = found.Text()
				}
			})

			valorBase := n.Find("li").First()
			priceString = valorBase.Text()

		})

	}

	log.Printf("SpaceString: %v\n", spaceString)
	log.Printf("RoomsString: %v\n", roomString)
	log.Printf("PriceString: %v\n", priceString)
	url := ctx.URL()
	log.Printf("URL: %v\n", url)

	titleParsed := parseTitle(title)
	location := parseLocationLifeImob(locationString)
	rooms, _ := parseRoomsLifeImob(roomString)
	space := parseSpaceLifeImob(spaceString)
	price := parsePriceLifeImob(priceString)
	log.Printf("Title: %v\n", titleParsed)
	log.Printf("Space: %v\n", space)
	log.Printf("Rooms: %v\n", rooms)
	log.Printf("Price: %v\n", price)
	log.Printf("URL: %v\n", url)

	item := Item{
		Title:       titleParsed,
		Location:    location,
		Price:       price,
		LivingSpace: space,
		Rooms:       rooms,
		Url:         url.String(),
		ScrapedAt:   time.Now().UTC(),
	}

	LocalExporter.storeDocument(item)

	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	if ctx.URL().Host == "www.lifeimob.com.br" {
		return true
	}
	return false
}
