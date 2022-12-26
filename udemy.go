package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
)

type Course struct {
	ID      string    `json:"-"`
	Created time.Time `json:"created"`
}

func main() {

	var url string

	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		fmt.Println("paste udemy url as first argument")
		os.Exit(0)
	}

	c := colly.NewCollector()
	var course Course
	course.ID = getCourseIDwithColly(url, c)

	resp, err := http.Get(fmt.Sprintf("https://www.udemy.com/api-2.0/courses/%s/?fields[course]=created", course.ID))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = json.Unmarshal(output, &course)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	color.Green("The course was created %v\n", course.Created.Format("Jan 2006"))
}

func getCourseIDwithColly(url string, c *colly.Collector) string {
	var courseID string

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("body", func(h *colly.HTMLElement) {
		courseID = h.Attr("data-clp-course-id")
	})

	c.Visit(url)

	return courseID
}
