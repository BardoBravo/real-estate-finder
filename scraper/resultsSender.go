package cloudrun

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"google.golang.org/api/iterator"
)

type ResultsSender interface {
	writeAndSend() error
}

type EmailSender struct {
	scraperName string
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func (emailer EmailSender) writeAndSend() {

	collectedRegistries := getNewRegistries(emailer.scraperName)

	if len(collectedRegistries) > 0 {
		sendEmail(collectedRegistries, emailer.scraperName)
	}

}

func getNewRegistries(scraperName string) string {
	currentTime := time.Now()
	timeLocation, errorTime := time.LoadLocation("America/Sao_Paulo")
	if errorTime != nil {
		log.Println("Error: ", errorTime)
	}

	log.Println("ScraperName: ", scraperName)

	searchedTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 0, 0, 0, 0, timeLocation)
	query := DBAccess.DBClient.Collection(scraperName).Where("ScrapedAt", ">", searchedTime).Documents(DBAccess.DBContext)
	collectedRegistries := ""
	counter := 0
	pairs := ""

	log.Printf("Query: %v", query)

	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		counter = counter + 1
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		m := doc.Data()
		for key, value := range m {
			log.Println(key, ": ", value)
			pairs = pairs + key + ": " + fmt.Sprintf("%v", value) + " | "
		}

		collectedRegistries = collectedRegistries + fmt.Sprintf("%v", counter) + " - " + pairs + "\n\n"
		pairs = ""
	}

	fmt.Println("New Registries: ", collectedRegistries)
	return collectedRegistries
}

func sendEmail(body string, scraper string) {
	from := "ImobiliariaCVTR@gmail.com"
	pass := "CVTRImob963"
	claudioEmail := "claudiovtramos@gmail.com"
	jaquelineEmail := "jackecassia@gmail.com"

	//TODO: change subject name so it contains the current date as well
	msg := "From: " + from + "\n" +
		"To: " + claudioEmail + "," + jaquelineEmail + "\n" +
		"Subject: Novos Imóveis" + scraper + "\n\n" + body

	log.Println("Msg: ", msg)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{jaquelineEmail, claudioEmail}, []byte(msg))

	if err != nil {
		log.Fatalf("smtp error: %s", err)
		return
	}

	log.Println("Email Sent")
}
