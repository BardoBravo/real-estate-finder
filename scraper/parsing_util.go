package cloudrun

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

//TODO: for version 1.3, check on how to format data into int / float, for better filtering in the CSV file
func parsePrice(valueStr string) (string, error) {
	replacer := strings.NewReplacer("R$ ", "", ",", "", ".", "")
	price, error := strconv.Atoi(replacer.Replace(valueStr))
	price = price / 100
	if error != nil {
		log.Print(error)
	}
	return valueStr, error
}

func parseSpace(value string) (float64, error) {
	return parseFloat(value, " m²")
}

func parseSpaceString(value string) string {
	replacer := strings.NewReplacer(" ", "")
	return replacer.Replace(value)
}

func parseFloat(valueStr string, unit string) (float64, error) {
	replacer := strings.NewReplacer(unit, "")
	return strconv.ParseFloat(replacer.Replace(valueStr), 64)
}

func parseURL(urlString string, space string) string {
	replacer := strings.NewReplacer(space, "")
	return replacer.Replace(urlString)
}

func parseTitle(value string) string {
	return strings.Replace(value, " C", "C", -1)
}

func parseRoomsLifeImob(value string) (float64, error) {
	log.Printf("To be Replaced: %v\n", value)
	replaced := strings.ReplaceAll(value, "Dormitório(s) ", "")
	log.Printf("Replaced: %v\n", replaced)
	return strconv.ParseFloat(replaced, 64)
}

func parseSpaceLifeImob(value string) string {
	replaced := strings.Replace(value, "Área Total ", "", -1)
	replaced = strings.Replace(replaced, "Área Terreno ", "", -1)
	replaced = strings.Replace(replaced, "Área privativa ", "", -1)
	return replaced
}

func parsePriceLifeImob(value string) string {
	regexNoTabs := regexp.MustCompile(`\x{0009}`)
	regexNoNewLines := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	replaced := regexNoTabs.ReplaceAllString(value, ``)
	replaced = regexNoNewLines.ReplaceAllString(replaced, ``)
	return strings.Replace(replaced, "Valor", "", -1)

}
