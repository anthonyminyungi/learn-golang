package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

func Scrape(term string) {
	var baseUrl string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	mch := make(chan []extractedJob)
	var jobs []extractedJob
	totalPages := getPages(baseUrl)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseUrl, mch)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-mch
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted: ", len(jobs))
}

func getPage(page int, baseUrl string, mch chan<- []extractedJob) {
	ch := make(chan extractedJob)
	var jobs []extractedJob
	pageURL := baseUrl + "&start=" + strconv.Itoa(50*page)
	fmt.Println("Requesting  ", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, ch)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-ch
		jobs = append(jobs, job)
	}

	mch <- jobs
}

func extractJob(card *goquery.Selection, ch chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find(".title>a").Text())
	location := CleanString(card.Find(".sjcl").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())
	ch <- extractedJob{id: id, title: title, location: location, salary: salary, summary: summary}
}

func getPages(baseUrl string) int {
	pages := 0
	res, err := http.Get(baseUrl)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func writeJobs(jobs []extractedJob) {
	// wch := make(chan bool)
	file, err := os.Create("jobs.csv")
	checkErr(err)
	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(utf8bom)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"link", "title", "location", "salary", "summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobLink := "https://kr.indeed.com/viewjob?jk=" + job.id
		jobSlice := []string{jobLink, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("request failed with status code:", res.StatusCode)
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
