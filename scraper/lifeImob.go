package cloudrun

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type LifeImob struct{}

func (LifeImob) parseItem(e *colly.XMLElement) Item {
	log.Println(e.ChildTexts("*"))
	selector := "//*[contains(@class, 'grid-itens')]//div(position = 2)//ul//"
	selector_previous := "//*[contains(contains(@id, 'divBusca')]//ul//li//"
	url := ""
	//url := e.ChildAttr(selector, "href")

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
	c := platform.NewCollector(config)

	c.OnXML("/html/body/div[5]/div/div/div[2]/div[2]/ul/li[1]", func(e *colly.XMLElement) {
		log.Print("Starting LifeImob Storage")
		item := platform.parseItem(e)

		//exporter.storeDocument(item)

		exporter.write(item)
	})

	/*c.OnXML("//div[contains(@class, 'pagination-bottom')]/div/a[last()]]", func(e *colly.XMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))
		log.Print("URL:" + url)
		url = parseURL(url, " ")
		log.Print("URL parsed:" + url)
		c.Visit(url)
	}) */

	c.Visit("https://www.lifeimob.com.br/index.asp?id_pagina=8&OrdernarPor=Preco%20DESC&ValorDe=0&ValorAte=500000000&cidade=3326")
	return c
}
