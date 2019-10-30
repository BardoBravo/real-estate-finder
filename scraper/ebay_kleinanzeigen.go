package cloudrun

import (
	"time"

	"github.com/gocolly/colly"
)

type EBayKleinanzeigen struct{}

func (EBayKleinanzeigen) parseItem(e *colly.XMLElement) Item {
	selector := "//*[contains(@class, 'ad-listitem')]//*[contains(@class, 'aditem-main')]//a"
	title := e.ChildText(selector)
	url := e.ChildAttr(selector, "href")
	locationNodes := e.ChildTexts("//*[contains(@class, 'aditem-details')]//text()")
	var location string
	if len(locationNodes) > 8 {
		location = locationNodes[6] + " " + locationNodes[8]
	}
	priceString := e.ChildText("//*[contains(@class, 'aditem-details')]//strong")
	price, _ := parsePrice(priceString)
	//spaceString := e.ChildText("//*[contains(@class, 'text-module-end')]//*[contains(text(), 'm²')]")
	//livingSpace, _ := parseSpace(spaceString)
	roomsString := e.ChildText("//*[contains(@class, 'text-module-end')]//*[contains(text(), 'Zimmer')]")
	rooms, _ := parseFloat(roomsString, " Zimmer")

	return Item{
		Title:            title,
		Location:         location,
		HasExactLocation: false,
		Price:            price,
		//livingSpace:      livingSpace,
		Rooms:     rooms,
		Url:       e.Request.AbsoluteURL(url),
		ScrapedAt: time.Now().UTC(),
	}
}

func (platform EBayKleinanzeigen) NewCollector(config Config) *colly.Collector {
	options := append(
		config.collectorOptions,
		colly.AllowedDomains("www.ebay-kleinanzeigen.de"))
	return colly.NewCollector(options...)
}

func (platform EBayKleinanzeigen) crawl(config Config, exporter Exporter) *colly.Collector {
	c := platform.NewCollector(config)

	c.OnXML("//*[contains(@class, 'ad-listitem')]", func(e *colly.XMLElement) {
		item := platform.parseItem(e)
		exporter.write(item)
	})

	c.OnXML("//a[contains(@class, 'pagination-next')]", func(e *colly.XMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(url)
	})

	c.Visit("https://www.ebay-kleinanzeigen.de/s-finder-imob-service-mieten/berlin/c203l3331")
	return c
}
