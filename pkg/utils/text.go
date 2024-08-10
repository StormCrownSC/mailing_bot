package utils

import (
	"strings"
	"time"
	"unicode"
)

func TextProcessing(text string) (answer []string) {
	formatText := strings.Fields(TextFormatting(text))
	for _, messageText := range formatText {
		var isID = true
		if len(messageText) < 8 {
			for _, letter := range []rune(messageText) {
				if !unicode.IsDigit(letter) {
					isID = false
				}
			}
		} else {
			for _, letter := range []rune(messageText) {
				if unicode.Is(unicode.Cyrillic, letter) {
					isID = false
				}
			}
		}

		if isID {
			answer = append(answer, messageText)
		}
	}
	return
}

func TextFormatting(text string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(text), "\u00a0", " "))
}

func parseTwoDates(commandText string) (time.Time, time.Time) {
	now := GetMoscowTime()

	dateStrings := strings.Fields(commandText)
	if len(dateStrings) == 0 {
		return now, now
	}

	if len(dateStrings) == 1 {
		layout := "2006-01-02"
		firstDate, err := time.Parse(layout, dateStrings[0])
		if err != nil {
			return now, now
		}
		return firstDate, firstDate
	}

	// Парсинг двух дат
	layout := "2006-01-02"
	firstDate, err := time.Parse(layout, dateStrings[0])
	if err != nil {
		return now, now
	}

	secondDate, err := time.Parse(layout, dateStrings[1])
	if err != nil {
		return firstDate, firstDate
	}

	return firstDate, secondDate
}

func FormatDate(commandText string) (string, string, string) {
	firstDate, secondDate := parseTwoDates(commandText)
	firstDateFormatted := firstDate.Format("2006-01-02")
	firstDateFormattedFormattedText := firstDate.Format("02.01.2006")

	secondDateFormatted := secondDate.Format("2006-01-02")
	secondDateFormattedFormattedText := secondDate.Format("02.01.2006")

	var textFormattedDate string
	if firstDateFormattedFormattedText != secondDateFormattedFormattedText {
		textFormattedDate = firstDateFormattedFormattedText + " - " + secondDateFormattedFormattedText
	} else {
		textFormattedDate = firstDateFormattedFormattedText
	}
	return firstDateFormatted, secondDateFormatted, textFormattedDate
}
