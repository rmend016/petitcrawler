package petitcrawler


import (
    "fmt"
    "flag"
    "strings"
    "os"
    "errors"
    "net/url"
)


// Set up custom flags from command line
var MaxpPtr= flag.Int( "maxprint", 10, "Maximum number of assests/children to print")
var UrlPtr = flag.String("url", "", "Starting URL to crawl. This is mandatory.")
var MaxcPtr = flag.Int("maxcrawl", 500, "Maximum number of pages to collect. Default 500.")
var MaxtPtr = flag.Int("maxtime", 60*3, "Max time in seconds to crawl for. Default 3 minutes.")
var HelpPtr = flag.Bool("help", false, "Help text." )
var OutfilePtr = flag.String("filename", "", "Specify a file to write the sitemap to. Default is <domain name>.txt .")
var NumwPtr = flag.Int("numworkers", 100, "The number of worker processes we spawn. Default is 100")

var MAX_WORKERS = 1000
 


func init() {
    flag.Parse()
    if *UrlPtr == "" || *HelpPtr == true {
        printHelp()
        flag.Usage()
        os.Exit(1)
    }
}


// Do a simple check on a URL - make sure we can extract the domain
func DomainCheck( domain *url.URL ) (error) {
    //domain, err := url.Parse(link)
    if domain.Host == "" {
        return errors.New("Unable to parse domain of URL.")
    }
    if strings.Contains(domain.Host, "www.") {
        domain.Host = strings.TrimLeft( domain.Host, "www." )
    }
    if domain.Scheme == "" {
        return errors.New("Must specify scheme in starting URL (ex: http).")
    }

    return nil
}



// printHelp prints information to the user when it was requested with cmd line flag help
func printHelp() {

    fmt.Println("\n---- Welcome to PetitCrawler! A single domain web crawler implemented in Golang! ----\n\n")
    fmt.Println("Example: ./Webcrawler -url http://www.urltocrawl.com\n")
    fmt.Println("This program was designed to crawl a single domain.\n")
    fmt.Println("The input to the program is a single URL in a format similar to: 'http://www.exampleurlnotreal.com'. Please follow this format as closely as possible to prevent any errors in crawling.\n")
    fmt.Println("There are some options that you can configure via command line, shown below. The program uses glog package to log any errors it encounters, exiting on fatal ones.\n")
    fmt.Println("The errors are very descriptive, and if you have an issue, you should be able to pinpoint what happened from the log.\n")
    fmt.Println("The output to the program is the site map of the single domain crawled, for each link crawled we display the: (1) URL, (2) static assets, (3) children links found on page.\n")
    fmt.Println("Thank you for using my program! If you have any suggestions for improvement, they are very welcome!")

}
