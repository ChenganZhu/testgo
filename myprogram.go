package main

import (
    "flag"
    "slices"
    "strings"
    "fmt"
    "golang.org/x/net/html"
    "net/http"
    "os"
    "encoding/json"
    "sort"
)

const output_stdout string = "stdout"
const output_json string = "json"

// Flag to handle a URL argument
type urlFlag []string
func (i *urlFlag) String() string {
    return ""
}
func (i *urlFlag) Set(value string) error {
    *i = append(*i, value)
    return nil
}

// Flag to handle the output format argument
type outputFlag struct {
    format string
}
func (i *outputFlag) String() string {
    return i.format
}
func (i *outputFlag) Set(value string) error {
    trimedValue := strings.ToLower(value) 
    outputs := []string{output_stdout, output_json}
    if slices.Contains(outputs, trimedValue) {
        i.format = trimedValue
    } else {
        return fmt.Errorf("must be either %s or %s", output_stdout, output_json)
    }
    return nil
}

func main() {
    var urls urlFlag
    var outputFormat outputFlag
    parseArguments(&urls,&outputFormat)

    var urlLinks = parseLinks(urls)

    displayLinks(urlLinks, &outputFormat)
}

/*
Function that defines and parses urls and output format.
*/
func parseArguments(urls *urlFlag, output *outputFlag) {
    flag.Var(urls, "u", "a URL target (can be specified multiple times)")
    flag.Var(output, "o", fmt.Sprintf("an output format, either %s or %s (default is %s)", output_stdout, output_json, output_stdout))
    flag.Parse()

    // default format is stdout
    if len(output.format) == 0 {
        output.format = output_stdout
    }

    if *urls == nil || len(*urls) == 0 {
        fmt.Println("u is not defined")
        flag.PrintDefaults()
        os.Exit(2)
    }
}

/*
Function that parses the provided urls and look for href element.
It returns a map where k is the provided url and v is a slice of links

This code was taken from https://www.makeuseof.com/parse-and-generate-html-in-go/
and slighty modified to handle edge cases such href with / witout root url and slashes.
*/
func parseLinks(urls []string) map[string][]string {
    urlLinks := make(map[string][]string)

    for _, url := range urls {
        // Send an HTTP GET request to the example.com web page
        // resp, err := http.Get("https://www.example.com")
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
        defer resp.Body.Close()

        // Use the html package to parse the response body from the request
        doc, err := html.Parse(resp.Body)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Find and print all links on the web page
        var links []string
        var link func(*html.Node)
        link = func(n *html.Node) {
            if n.Type == html.ElementNode && n.Data == "a" {
                for _, a := range n.Attr {
                    // adds a new link entry when the attribute matches
                    if a.Key == "href" {
                        // when child link, prefix with slash if none
                        var sb strings.Builder
                        if hasProtocol(a.Val) || strings.HasPrefix(a.Val, "/") {
                            sb.WriteString(a.Val)
                        } else {
                            sb.WriteString("/")
                            sb.WriteString(a.Val)
                        }
                        links = append(links, sb.String())
                    }
                }
            }

            // traverses the HTML of the webpage from the first child node
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                link(c)
            }
        }
        link(doc)

        sort.Strings(links)
        urlLinks[url] = links
    }
    return urlLinks
}

/*
Function that output the parsed links, either as json or standard output
*/
func displayLinks(urlLinks map[string][]string, output *outputFlag) {
    if output.format == output_json {
        jsonStr, err := json.Marshal(urlLinks)
        if err != nil {
            fmt.Printf("Error: %s", err.Error())
        } else {
            fmt.Println(string(jsonStr))
        }
    } else {
        // stdout display
        for root, links := range urlLinks {
            for _, link := range links {
                if hasProtocol(link) {
                    // no need to display root url
                    fmt.Printf("%s\n\n", link)
                } else {
                    // concat root url with link
                    fmt.Printf("%s%s\n\n", root, link)
                }
            }
        }
    }
}

/*
Check if the string starts with protocols http or https
*/
func hasProtocol(str string) bool {
    return strings.HasPrefix(str,"http") || strings.HasPrefix(str,"https")
}
