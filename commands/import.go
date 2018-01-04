package commands

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Mentatious/mentat-cli/io"
	"github.com/Mentatious/mentat-cli/schema"
	"github.com/ybbus/jsonrpc"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ImportDelicious ... import file with Delicious exported bookmarks
func ImportDelicious(filename string, rpcClient *jsonrpc.RPCClient, userID string, quiet bool) {
	log := io.GetLog()
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s, exiting...", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	var posts schema.Posts
	b, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(b, &posts)

	countWhole := 0
	countFailed := 0
	var failedData []schema.ImportError

	timeStart := time.Now()
	log.Infof("Found %d bookmarks", posts.Total)

	for _, post := range posts.PostsList {
		countWhole++
		link := post.Href
		tags := strings.Split(post.Tags, " ")
		metadata := schema.PostMetadata{
			Description:     post.Description,
			TimeAddedOrigin: post.Time,
			HashOrigin:      post.Hash,
			MetaOrigin:      post.Meta,
			From:            "delicious",
		}
		rpcResp, err := rpcClient.CallNamed("entry.Add", map[string]interface{}{
			"userid":   userID,
			"type":     "bookmark",
			"content":  link,
			"tags":     tags,
			"metadata": metadata,
		})
		if err != nil {
			countFailed++
			failedData = append(failedData, schema.ImportError{
				Link:  link,
				Error: err.Error(),
			})
		} else {
			respData := schema.AddResponse{}
			err := rpcResp.GetObject(&respData)
			if err != nil {
				log.Infof("link: %s; error: %s", link, err)
			} else {
				log.Infof("link: %s; result: %s", link, respData.Message)
			}
		}
	}
	timeFinish := time.Now()
	log.Infof("Processed %d bookmarks in %v", posts.Total, timeFinish.Sub(timeStart))
	if countFailed > 0 {
		if quiet {
			log.Infof("%d of them failed", countFailed)
		} else {
			log.Infof("%d of them failed, details below:", countFailed)
			for _, fail := range failedData {
				log.Infof("link: %s; error: %s", fail.Link, fail.Error)
			}
		}
	}
}

// ImportPocket ... import file with Pocket exported bookmarks
func ImportPocket(filename string, rpcClient *jsonrpc.RPCClient, userID string, quiet bool) {
	log := io.GetLog()
	htmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s, exiting...", err)
		os.Exit(1)
	}
	defer htmlFile.Close()
	contents, _ := ioutil.ReadAll(htmlFile)

	root, err := html.Parse(bytes.NewReader(contents))
	if err != nil {
		panic(err)
	}

	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent.DataAtom == atom.Li && n.Parent.Parent.DataAtom == atom.Ul {
			return true
		}
		return false
	}

	bmarks := scrape.FindAll(root, matcher)
	bmarksCount := len(bmarks)
	countWhole := 0
	countFailed := 0
	var failedData []schema.ImportError
	timeStart := time.Now()
	for _, bmark := range bmarks {
		link := scrape.Attr(bmark, "href")
		tags := strings.Split(scrape.Attr(bmark, "tags"), " ")
		countWhole++
		metadata := schema.PostMetadata{
			Description:     scrape.Text(bmark),
			TimeAddedOrigin: scrape.Attr(bmark, "time_added"),
			From:            "pocket",
		}
		rpcResp, err := rpcClient.CallNamed("entry.Add", map[string]interface{}{
			"userid":   userID,
			"type":     "bookmark",
			"content":  link,
			"tags":     tags,
			"metadata": metadata,
		})
		if err != nil {
			countFailed++
			failedData = append(failedData, schema.ImportError{
				Link:  link,
				Error: err.Error(),
			})
		} else {
			respData := schema.AddResponse{}
			err := rpcResp.GetObject(&respData)
			if err != nil {
				log.Infof("link: %s; error: %s", link, err)
			} else {
				log.Infof("link: %s; result: %s", link, respData.Message)
			}
		}
	}
	timeFinish := time.Now()
	log.Infof("Processed %d bookmarks in %v", bmarksCount, timeFinish.Sub(timeStart))
	if countFailed > 0 {
		if quiet {
			log.Infof("%d of them failed", countFailed)
		} else {
			log.Infof("%d of them failed, details below:", countFailed)
			for _, fail := range failedData {
				log.Infof("link: %s; error: %s", fail.Link, fail.Error)
			}
		}
	}
}

// ImportOrg ... import Org-mode file
func ImportOrg(filename string, rpcClient *jsonrpc.RPCClient) {
}

// import from Facebook saved links (first export to html)
// #!/usr/bin/env ruby

// require 'nokogiri'
// require 'uri'
// require 'cgi'

// input_file = ARGV[0]

// page = Nokogiri.HTML(File.read(input_file))
// page.xpath("//div[@class='_4bl9 _5yjp']").each() do |d|
// 	raw_link = d.xpath("a")[0].attributes["href"].value
//     link = URI(raw_link)
//     query_string = link.query
//     if query_string.nil?
//     puts link
//     else
//     puts CGI.parse(query_string)["u"]
//     end
// end
