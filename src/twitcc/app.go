package main

import (
	//built-in libs
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	removedNewLine chan []string
	trigger        bool
	timer          *time.Timer
	twitterAcct    string
	sec            string
	durtion        time.Duration
)

//This function used to handle err biolerplates
func checkErr(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

//This function will display/print any value passed to it
//for dubugging purposes, that's why interface{} is used
func debug(v interface{}) {
	fmt.Printf("%v\n", v)
}

func init() {
	removedNewLine = make(chan []string, 1)
	flag.StringVar(&twitterAcct, "a", "", "Twitter Account Name")
	flag.StringVar(&sec, "t", "3s", "Set Timer Peroid")
	flag.Parse()
	durtion, err := time.ParseDuration(sec)
	checkErr(err)
	timer = time.NewTimer(durtion)
}

//This function will fetch all tweets from specified
//twitter account and after that they will be passed
// to []string channel for further implementation in
// main funuction
func fetch() {

	resp, err := http.Get("https://twitter.com/" + twitterAcct)
	checkErr(err)
	content, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	defer resp.Body.Close()
	//This doesn't work with pinnned tweets must be removed from timeline
	removedNewLine <- strings.Split(string(content), "\n")

}

func main() {
	for {
		if twitterAcct != "" {
			go fetch()

			select {
			case r := <-removedNewLine:
				if trigger {
					for i := 0; i < len(r); i++ {
						if strings.Contains(r[i], "data-aria-label-part=\"0\">") {
							command := strings.Split(strings.Split(r[i], "data-aria-label-part=\"0\">")[1], "</p>")[0]
							debug(command)
							//return
							i = len(r)
						}
					}
					trigger = false
				}
			case <-timer.C:
				//every 3 sec tick
				trigger = true
				//debug("Timer Triggered")
				timer.Reset(durtion)
			}
		} else {
			debug("Please make sure to use the right option!")
			return
		}
	}
}
