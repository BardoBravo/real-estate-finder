package cloudrun

import (
	"fmt"
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

func (LifeImob) parseItem(e *colly.HTMLElement) Item {
	//LocalExporter.storeDocument()
	log.Print("teste")
	log.Print(e)
	/* log.Println(e.ChildTexts("*"))
	log.Println(e.ChildAttrs("*", "h1"))

	iTags := e.DOM.Find("i")
	iTags.Each(func(_ int, s *goquery.Selection) {
		log.Println(s)
	})
	log.Println(e.DOM.Find("font").Text()) */

	selector := "div[class=grid-list-header]"
	log.Println(e.ChildTexts(selector))
	selector_previous := "//*[contains(contains(@id, 'divBusca')]//ul//li//"
	//url := ""
	url := e.ChildAttr(selector, "href")

	title := e.ChildText("//a//div[position() = 3//div[first()//p//text()]]")

	location := e.ChildText("//*[contains(contains(@id, 'divBusca')]//ul//li//a//div[position() = 3//h1//text()]]")
	location = location + " em " + e.ChildText("//*[contains(contains(@id, 'divBusca')]//ul//li//a//div[position() = 3//h2//text()]]")
	priceString := e.ChildText("//*[contains(contains(@id, 'divBusca')]//ul//li//a//div[position() = 2//span//text()]]")
	price, _ := parsePrice(priceString)
	spaceString := e.ChildText("//*[contains(@class, 'pull-left')]//span[position() = 4]//text()")
	livingSpace := parseSpaceString(spaceString)
	roomsString := e.ChildText("//*[contains(@class, 'pull-left')]//span[position()= 1]//text()]")
	rooms, _ := parseFloat(roomsString, " quartos")
	log.Print("-> Selector:" + selector +
		"-> Title: " + title +
		"-> URL: " + url +
		"-> Selector Previous: " + selector_previous +
		"-> Price: " + price +
		"-> PriceString: " + priceString +
		"-> SpaceSpring: " + spaceString +
		"-> livingSpace: " + livingSpace +
		"-> roomString: " + roomsString +
		"-> rooms: " + fmt.Sprintf("%.2f", rooms))

	return Item{
		Title:       title,
		Location:    location,
		Price:       price,
		LivingSpace: livingSpace,
		Rooms:       rooms,
		Url:         e.Request.AbsoluteURL(url),
		ScrapedAt:   time.Now().UTC(),
	}
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
	opts.MaxVisits = 30

	log.Print("starting crawl...")
	c := gocrawl.NewCrawlerWithOptions(opts)
	if err := c.Run("https://www.lifeimob.com.br/index.asp?id_pagina=8&OrdernarPor=Preco%20DESC&ValorDe=0&ValorAte=500000000&cidade=3326"); err != nil {
		log.Print(err)
	}

	//c.Visit("https://www.lifeimob.com.br/index.asp?id_pagina=8&OrdernarPor=Preco%20DESC&ValorDe=0&ValorAte=500000000&cidade=3326")
	return colly
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())

	spaceString := ""
	priceString := ""
	roomString := ""
	title := ""
	location := ""

	star := doc.Find(".texts-center")
	if star.Length() > 0 {
		star.Each(func(i int, n *goquery.Selection) {
			link := n.Find("b")
			title = link.Text()
			log.Printf("Title: %v\n", title)

			location = n.Find("h1").Text()
			log.Printf("Location: %v\n", location)

		})
	}

	data := doc.Find(".icons-dd")
	if data.Length() > 0 {
		data.Each(func(i int, n *goquery.Selection) {
			n.Find("li").Each(func(cont int, found *goquery.Selection) {
				dadosGenericos := found.Text()
				log.Printf("All: %v\n", dadosGenericos)

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
