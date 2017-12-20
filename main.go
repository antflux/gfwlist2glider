package main

// https://code.google.com/p/autoproxy-gfwlist/wiki/Rules
// crontab -e
// 9 5 * * * cd /opt/gfwlist2glider && /opt/gfwlist2glider/gfwlist2glider > /dev/null 2>&1

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	var GfwlistURL = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"
	// var GfwlistURL = "https://raw.githubusercontent.com/gfwlist/tinylist/master/tinylist.txt"
	var GliderListFile = "gfwlist.list"
	var RegexComment = `^\!|\[|^@@|^\d+\.\d+\.\d+\.\d+`
	var RegexDomain = `([\w\-\_]+\.[\w\.\-\_]+)[\/\*]*`

	log.Println("fetching gfwlist...")
	resp, err := http.Get(GfwlistURL)
	if err != nil {
		log.Fatal("error:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error:", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Fatal("error:", err)
	}

	GliderList, err := os.Create(GliderListFile)
	if err != nil {
		log.Fatal("error:", err)
	}
	defer GliderList.Close()

	gfwlist := string(data)

	lines := strings.Split(gfwlist, "\n")
	reComment := regexp.MustCompile(RegexComment)
	reDomain := regexp.MustCompile(RegexDomain)

	GliderList.WriteString("# gfwlist for glider\n")
	GliderList.WriteString("# updated on " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	GliderList.WriteString("#\n")

	domains := make(map[string]bool)

	//GFWLists
	for _, line := range lines {
		if submatch := reComment.FindString(line); submatch != "" {
			fmt.Println("COMMENT LINE: ", line)
		} else if submatch := reDomain.FindString(line); submatch != "" {
			domain := strings.TrimSuffix(submatch, "*")
			domain = strings.TrimSuffix(domain, "/")
			fmt.Println("DOMAIN LINE: ", domain, line)

			domains[domain] = true

		}
	}

	for domain := range domains {
		if ip := net.ParseIP(domain); ip != nil {
			GliderList.WriteString("ip=")
		} else {
			GliderList.WriteString("domain=")
		}

		GliderList.WriteString(domain)
		GliderList.WriteString("\n")
	}

}
