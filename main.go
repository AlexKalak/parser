package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const url string = "http://iq.cq.md/"

var file *os.File

func main() {
	t := time.Now()
	var wg sync.WaitGroup
	var err error
	file, err = os.OpenFile("./log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, fs.ModeAppend)

	if err != nil {
		panic(err)
	}

	threads := 8
	start := 800_000
	end := 810_000
	change := (end - start) / threads

	for i := 0; i < threads; i++ {
		startIndex := start + change*i
		endIndex := startIndex + change

		fmt.Println(startIndex)
		fmt.Println(endIndex)

		wg.Add(1)
		go doQueries(startIndex, endIndex, &wg)
	}

	wg.Wait()
	fmt.Printf("Time: %v", time.Since(t).Seconds())
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}

func doQueries(start int, end int, wg *sync.WaitGroup) {
	for i := start; i < end; i++ {
		// fmt.Println(i)
		currentUrl := url + strconv.Itoa(i)

		resp, err := http.Get(currentUrl)
		if err != nil {
			panic(err)
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		doc := string(bytes)

		title := getTitle(&doc)

		if strings.Contains(title, "ЮГ") {
			fmt.Println(title)
			_, err := file.WriteString(currentUrl + "	:			" + title + "\n\n")
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("DONE")
	(*wg).Done()
}

func getTitle(doc *string) string {
	start := strings.Index(*doc, "<title>")
	if start == -1 {
		title := ""
		return title
	}
	start += len("<title>")
	end := strings.Index(*doc, "</title>")
	title := (*doc)[start:end]
	return title
}
