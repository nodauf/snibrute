package brute

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Options struct {
	Domain       string
	Wordlist     string
	Ip           string
	MatchStatus  string
	MatchSize    string
	FilterStatus string
	FilterSize   string
	Port         int
	Verbose      bool
}

func (option *Options) Brute() {
	subdomains := readFile(option.Wordlist)
	for _, subdomain := range strings.Split(subdomains, "\n") {
		if subdomain == "" {
			continue
		}
		host := subdomain + "." + option.Domain

		req, _ := http.NewRequest("GET", "https://"+option.Ip+":"+strconv.Itoa(option.Port)+"/", nil)

		req.Host = host

		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					ServerName:         req.Host,
					InsecureSkipVerify: true,
				},
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			if option.Verbose {
				fmt.Println("For host " + color.RedString(host) + " error: " + err.Error())
			}
		} else if option.match(resp) && !option.filter(resp) {
			dataB, _ := ioutil.ReadAll(resp.Body)
			data := string(dataB)

			fmt.Println(color.GreenString(subdomain) + "\t\t\t" + "[Status: " + strconv.Itoa(resp.StatusCode) + ", Size: " + strconv.Itoa(int(resp.ContentLength)) + ", Lines: " + strconv.Itoa(len(strings.Split(data, "\n"))) + "]")
		}

	}
}

// filter return true if the response is filtered and should not be displayed
func (option *Options) filter(resp *http.Response) bool {
	statusCodeSlice := strings.Split(option.FilterStatus, ",")
	sizeSlice := strings.Split(option.FilterSize, ",")
	return stringInSlice(statusCodeSlice, strconv.Itoa(resp.StatusCode)) || stringInSlice(sizeSlice, strconv.Itoa(int(resp.ContentLength)))
}

// match return true if the response matched and should be displayed
func (option *Options) match(resp *http.Response) bool {
	matchStatus := false
	matchSize := false
	if option.MatchStatus != "" {
		statusCodeSlice := strings.Split(option.MatchStatus, ",")
		matchStatus = stringInSlice(statusCodeSlice, strconv.Itoa(resp.StatusCode))
	}

	if option.MatchSize != "" {
		sizeSlice := strings.Split(option.MatchSize, ",")
		matchSize = stringInSlice(sizeSlice, strconv.Itoa(int(resp.ContentLength)))
	}
	// If no match configuration was asked return all results
	if option.MatchStatus == "" && option.MatchSize == "" {
		return true
	}
	return matchStatus || matchSize
}

func readFile(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error while reading the file " + file + ":" + err.Error())
		os.Exit(1)
	}
	if len(data) == 0 {
		fmt.Println("File " + file + " is empty")
		os.Exit(1)
	}
	return string(data)
}

func stringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
