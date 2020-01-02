package cloudrun

import (
	"log"
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
