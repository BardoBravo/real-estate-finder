package cloudrun

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Jus_sl struct{}

func (Jus_sl) parseItem(e *colly.XMLElement) Item {
	selector := "//*[contains(@class, 'col-sm-6 col-md-4 boxi animated')]//div//a"
	selector_previous := "//div/div"
	title := e.ChildAttr(selector_previous, "class")
	url := e.ChildAttr(selector, "href")
	location := e.ChildText("//*[contains(@class, 'onde')]//text()")
	location = location + " " + e.ChildText("//*[contains(@class, 'cat')]//text()")
	priceString := e.ChildText("//*[contains(@class, 'valor')]//strong")
	price, _ := parsePrice(priceString)
	spaceString := e.ChildText("//*[contains(@class, 'pull-right')]//span[position() = 4]//text()")
	livingSpace := parseSpaceString(spaceString)
	roomsString := e.ChildText("//*[contains(@class, 'pull-right')]//span[position()= 1]//text()]")
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

func (platform Jus_sl) NewCollector(config Config) *colly.Collector {
	log.Print("Starting Justus2")
	options := append(
		config.collectorOptions,
		colly.AllowedDomains("www.justoimoveis.com.br"))
	return colly.NewCollector(options...)
}

func (platform Jus_sl) crawl(config Config, exporter Exporter) *colly.Collector {
	c := platform.NewCollector(config)

	c.OnXML("//*[contains(@class, 'col-sm-6 col-md-4 boxi animated')]", func(e *colly.XMLElement) {

		log.Print("Starting Justus1")
		item := platform.parseItem(e)

		exporter.storeDocument(item)

		exporter.write(item)
	})

	c.OnXML("//div[contains(@class, 'paginacao-bottom')]/a[last()]]", func(e *colly.XMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))
		log.Print("URL:" + url)
		url = parseURL(url, " ")
		log.Print("URL parsed:" + url)
		c.Visit(url)
	})

	c.Visit("http://www.justoimoveis.com.br/imoveis?busca=venda&finalidade=venda&cidade=sao-leopoldo")
	return c
}
