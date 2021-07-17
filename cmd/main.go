package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ip-rw/req/pkg/client"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

// variables. Capital means public scope lowerCase means hidden from other packages
var (
	workers = flag.Int("workers", 10, "number of workers")
	verbose = flag.Bool("verbose", false, "verbose")

	titlePat = regexp.MustCompile("(?:<title>)(.+)(?:</title>)")
)

type Response struct {
	URL         *url.URL
	Code        int
	Size        int
	ContentType string
	Title       string
	Err         error // Mayybee
}

// called BEFORE main()
func init() {
	flag.Parse()
	log.SetLevel(log.InfoLevel)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func stdinProducer(stdin chan string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		stdin <- strings.TrimSpace(s.Text())
	}
	close(stdin)
}

func worker(stdin chan string, resultChan chan *Response) {
	for uri := range stdin {
		res, err := getUrl(uri)
		if err != nil {
			log.WithError(err).Debugf("%s failed", uri)
			continue
		}
		resultChan <- res
	}
}

func output(resultChan chan *Response) {
	for res := range resultChan {
		log.WithFields(log.Fields{
			"URL":         res.URL,
			"Code":        res.Code,
			"ContentType": res.ContentType,
			"Size":        res.Size,
			"Title":       res.Title,
		}).Info("success")
	}
}

// entry point
func main() {
	wg := &sync.WaitGroup{}
	stdinChan := make(chan string, *workers*2)
	resultChan := make(chan *Response, *workers*2)
	// lets  get  starrted with stdin. show yoou a litle tricck  ;)

	wg.Add(1)
	go func() {
		stdinProducer(stdinChan)
		wg.Done()
	}()

	wg.Add(*workers) // thhink  abbout itt, why  not?
	for i := 0; i < *workers; i++ {
		go func() {
			worker(stdinChan, resultChan)
			wg.Done()
		}()
	}

	go output(resultChan)
	wg.Wait()
	close(resultChan)
}

func getUrl(uri string) (*Response, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	c, err := client.NewHostClient(fmt.Sprintf("%s:%s", u.Hostname(), port), "", u.Scheme == "https")
	if err != nil {
		log.Debugln(err)
		return nil, err
	}
	defer c.Release()
	//fmt.Println(u, u.String())
	code, body, err := c.Get(nil, u.String())
	if err != nil {
		log.Warnln(err)
		return nil, err
	}

	mime := http.DetectContentType(body)
	ti := titlePat.FindSubmatch(body)
	title := ""
	if ti != nil {
		title = string(ti[1])
	}

	return &Response{
		URL:         u,
		Code:        code,
		Size:        len(body),
		ContentType: mime,
		Title:       title,
		//Err:         err,
	}, err
}


