package cloudrun

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Rodrigues struct{}

func (Rodrigues) parseItem(e *colly.XMLElement) Item {
	selector := "//*[contains(@class, 'col-sm-12 col-lg-6 box-align')]//div//a"
	selector_previous := "//div//a"
	title := e.ChildText("//*[contains(@class, 'card-text')]//text()")
	url := e.ChildAttr(selector, "href")
	location := e.ChildText("//*[contains(@class, 'card-title')]//text()")
	priceString := e.ChildText("//*[contains(@class, 'money location')]//text()")
	price, _ := parsePrice(priceString)
	spaceString := e.ChildText("//*[contains(@class, 'values')]//div[position() = 5]//p//span//text()")
	livingSpace := parseSpaceString(spaceString)
	roomsString := e.ChildText("//*[contains(@class, 'values')]//div[position() = 1]//p//span//text()")
	rooms, _ := parseFloat(roomsString, " quartos")
	log.Print("-> Selector:" + selector +
		"-> Title: " + title +
		"-> Selector Previous: " + selector_previous +
		"-> Price: " + price +
		"-> PriceString: " + priceString +
		"-> SpaceSpring: " + spaceString +
		"-> livingSpace: " + livingSpace +
		"-> roomString: " + roomsString +
		"-> rooms: " + fmt.Sprintf("%.2f", rooms))

	return Item{
		Title:    title,
		Location: location,
		//HasExactLocation: false,
		Price:       price,
		LivingSpace: livingSpace,
		Rooms:       rooms,
		Url:         e.Request.AbsoluteURL(url),
		ScrapedAt:   time.Now().UTC(),
	}
}

func (platform Rodrigues) NewCollector(config Config) *colly.Collector {
	log.Print("Starting Rodrigues")
	options := append(
		config.collectorOptions,
		colly.AllowedDomains("www.rodrigues.imb.br"))
	return colly.NewCollector(options...)
}

func (platform Rodrigues) crawl(config Config, exporter Exporter) *colly.Collector {
	c := platform.NewCollector(config)

	c.OnXML("//*[contains(@class, 'col-sm-12 col-lg-6 box-align')]", func(e *colly.XMLElement) {

		log.Print("Starting Rodrigues Card")
		item := platform.parseItem(e)
		exporter.write(item)
	})

	/* c.OnXML("//div[contains(@class, 'pagination-cell')]/a", func(e *colly.XMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))
		log.Print("URL:" + url)
		url = parseURL(url, " ")
		log.Print("URL parsed:" + url)
		c.Visit(url)
	}) */

	c.Visit("http://www.rodrigues.imb.br/imoveis/a-venda/sao-leopoldo")
	return c
}
