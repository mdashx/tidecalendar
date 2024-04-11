package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func QueryAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}

type tideTime struct {
	tide string
	when time.Time
}

func fetchTideTimes(monthStr string) *html.Node {
	cachePath := ".cache"
	err := os.MkdirAll(cachePath, 0750)

	docName := fmt.Sprintf("puntarenas-calendar-%s.htm", monthStr)
	docPath := filepath.Join(cachePath, docName)

	var doc *html.Node

	contents, err := os.ReadFile(docPath)
	if errors.Is(err, os.ErrNotExist) {
		url := fmt.Sprintf("https://www.tidetime.org/central-america/costa-rica/%s", docName)
		resp, err := http.Get(url)
		check(err)
		defer resp.Body.Close()

		// https://stackoverflow.com/questions/9644139/from-io-reader-to-string-in-go
		buf := new(strings.Builder)
		_, err = io.Copy(buf, resp.Body)
		check(err)
		os.WriteFile(docPath, []byte(buf.String()), 0644)
		log.Print("fetching HTML")
		doc, err = html.Parse(strings.NewReader(buf.String()))
		check(err)
	} else {
		doc, err = html.Parse(bytes.NewReader(contents))
	}

	return doc
}

func parseTideTimes(month string, monthStr string) []tideTime {
	doc := fetchTideTimes(monthStr)
	var headers []string
	var tideTimes []tideTime
	year := "2024"

	for i, r := range QueryAll(doc, "tr") {
		if i == 0 {
			for _, h := range QueryAll(r, "th") {
				if h.FirstChild != nil {
					headers = append(headers, h.FirstChild.Data)
				}
			}
		}

		// Skipping the first row because it is headers
		if i > 0 {
			var eventDate string
			for _, h := range QueryAll(r, "th") {
				if h.FirstChild != nil {
					day := strings.Split(h.FirstChild.Data, " ")
					dayNum := day[len(day)-1]
					eventDate = fmt.Sprintf("%s-%s-%s", year, month, dayNum)
				}
			}

			for j, d := range QueryAll(r, "td") {
				if d.FirstChild != nil {
					// Add 1 because the first cell in the row is a <th> element
					// So our first header column in the slice is extraneous
					h := headers[j+1]
					if h == "High" || h == "Low" {
						t := eventDate + " " + strings.SplitAfter(d.FirstChild.Data, "CST")[0]
						p, err := time.Parse("2006-01-02 3:04 PM MST", t)
						check(err)
						tideTimes = append(tideTimes, tideTime{tide: h, when: p})
					}
				}
			}
		}
	}
	return tideTimes
}

func createCalendar(tideTimes []tideTime, whichTides []string) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)

	color := make(map[string]string)
	color["High"] = "#47b0ff"
	color["Low"] = "#bad700"

	for _, t := range tideTimes {
		if slices.Contains(whichTides, t.tide) {
			event := cal.AddEvent(uuid.NewString())
			event.SetDtStampTime(time.Now())
			event.SetSummary(fmt.Sprintf("%s Tide", t.tide))
			event.SetStartAt(t.when)
			event.SetEndAt(t.when.Add(time.Minute * 15))
			event.SetColor(color[t.tide])
			event.SetDescription("Visit www.mdashx.com/tides to download or subscribe this calendar. Check back in 2025 for the updated 2025 calendar.")
		}
	}

	docName := "tides.ics"
	os.WriteFile(docName, []byte(cal.Serialize()), 0644)
}

func main() {
	months := [][]string{{"04", "apr"}, {"05", "may"}, {"06", "jun"}, {"07", "jul"}, {"08", "aug"}, {"09", "sep"}, {"10", "oct"}, {"11", "nov"}, {"12", "dec"}}
	tides := []tideTime{}

	for _, m := range months {
		tides = append(tides, parseTideTimes(m[0], m[1])...)
	}

	// whichTides := []string{"Low"}
	whichTides := []string{"High"}
	// whichTides := []string{"Low", "High"}
	createCalendar(tides, whichTides)
}
