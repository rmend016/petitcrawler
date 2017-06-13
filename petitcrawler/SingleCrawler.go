package petitcrawler

import (
    "fmt"
    "net/url"
    "os"
    "errors"
    "time"
    "strings"
    "sync"
    "github.com/asaskevich/govalidator"
    "github.com/golang/glog"
)


// A struct to represent a web crawler that will only crawl the domain it's initialized with
type SingleCrawler struct{

    Site *url.URL           // single site/ domain to be crawled
    Sitemap [] Page         // a sitemap made of Pages
    NumPages int            // number of pages collected - that are unique
    NumWorkers int          // number of workers to spawn 
    PRINT_LIMIT int         // for printing the site map, only display this many assets
    MAX_PAGES int           // max pages to crawl
    MAX_TIME time.Duration  // max time to crawl
    Filename string         // option to output sitemap to a file

}


// IsOk checks if a given crawler is okay to use 
func IsOk(c *SingleCrawler) error {
    if c.Site == nil {
        return errors.New("Crawler has no Site.")
    }
    if c.Sitemap == nil {
        return errors.New("Crawler has no Sitemap.")
    }
    if c.NumPages < 0 {
        return errors.New("Crawler has negative # of pages.")
    }
    if c.NumWorkers <= 0 {
        return errors.New("Crawler <= 0 number of workers (can't work).")
    }
    if c.Filename == "" {
        return errors.New("Crawler has no Filename to write sitemap to.")
    }
    return nil
}


// NewCrawler creates a new SingleCrawler instance, initializing all fields, 
// given a starting URL. 
func New() (*SingleCrawler, error) {

    defer glog.Flush()

    var crawler SingleCrawler
    startURL := *UrlPtr
    maxp := *MaxpPtr
    maxc := *MaxcPtr
    maxt := *MaxtPtr
    Filename := *OutfilePtr
    NumWorkers := *NumwPtr

    // validate the user input URL and decide if it's okay to use
    if govalidator.IsURL(startURL) == false {
        glog.Error("The starting URL is invalid. Please enter a valid URL.")
        return nil, errors.New("Bad starting URL.")
    }
    if maxp < 0 || maxc < 0 || maxt < 0 {
        glog.Error("Please pass in values > = 0 for max constraints (max print, max pages, max time). Please pass > 0 for the number of workers.")
        return nil, errors.New("Bad values for maxprint, maxpages, maxtime, or NumWorkers")
    }
    if NumWorkers <= 0 || NumWorkers > MAX_WORKERS {
        glog.Error("Number of workes is invalid. Must be > 0, and less that MAX_WORKERS.")
        return nil, errors.New("Bad value for NumWorkers")
    }
    if len(Filename) >= 255 {
        glog.Error("Filename can't be larger than 255 characters. Trimming Filename.")
        Filename = Filename[0:100]
    }


    crawler.MAX_PAGES = maxc
    crawler.PRINT_LIMIT = maxp
    crawler.NumPages = 0
    crawler.NumWorkers = NumWorkers
    crawler.MAX_TIME = time.Duration(maxt) * time.Second
    crawler.Sitemap = make( [] Page, crawler.MAX_PAGES)
    

    // Parse the URL - make sure it's ok to use
    domain, err := url.Parse(startURL)
    if err != nil {
        glog.Error("Error parsing domain of starting URL")
        return nil, errors.New("Unable to parse domain of start URL.")
    }
    err = DomainCheck( domain )
    if err != nil {
        glog.Error("Error parsing domain of starting URL")
        return nil, err
    }
    crawler.Site = domain
    
    if Filename != "" {
        crawler.Filename = Filename
    } else {
        crawler.Filename = crawler.Site.Host + ".txt"
        if len( crawler.Filename ) >= 255 {
            crawler.Filename = crawler.Filename[0:100]
        }
    }

    if err = IsOk( &crawler ); err!=nil{
        return nil, err
    }

    return &crawler, nil

}


// Start begins the crawling process based on the starting url 
// by passing urls in it's URL list to worker threads
func ( crawler *SingleCrawler ) Start()(error) {

    defer glog.Flush()

    if err1 := IsOk( crawler ); err1!=nil{
        return err1
    }

    // Stats for termination conditions 
    t0 := time.Now()            //Terminate after a given time
    noIncrease := 0             //Make sure we are finding unique sites, and not in an inf loop 
    last_pagecount := 0         //Keep track of the last # of pages
    var wg sync.WaitGroup           //For termination, to wait on workers

    // Channels for communication to workers
    pages := make( chan Page, crawler.NumWorkers*10 )
    rurls := make( chan string, crawler.NumWorkers*10 )
    surls := make( chan string, crawler.NumWorkers*10 )
    shutdown := make( chan bool, crawler.NumWorkers )

    // Map for making pages and urls unique
    assets := make( map[string][]string )
    vList := make( map[string]int )
 
    // Start the crawling, by providing the inital site URL
    surls <- crawler.Site.String()
    vList[crawler.Site.String()]++
    

    // Spawn the requested number of workers for the program
    for i:= 0; i< crawler.NumWorkers; i++ {
        wg.Add(1)
        go Worker( i, surls, rurls, crawler.Site, pages, shutdown, &wg )
    }


    for {

        select { 

            case link := <- rurls:
                // Receive a link to crawl, make sure it's unvisited, then send back
                if _, ok := vList[link]; ok == false {
                    glog.Info( fmt.Sprintf("starting crawler for %s\n", link))
                    surls <- link
                    vList[link]++
                } 

            case p := <- pages:
                //receive a page in the page channel, append it to the crawler's sitemap, if it's unique.
                ind := strings.Join(p.Assets, " ")
                if crawler.NumPages < len(crawler.Sitemap){
                    if _, ok := assets[ind]; ok == false {
                        assets[ind] = p.BabyUrls
                        crawler.Sitemap[crawler.NumPages] = p
                        crawler.NumPages += 1
                    }
                }
            default:
                // Print status update. 
                // Check termination conditions: time, space, nonincreasing, no urls left
                if time.Since(t0) % 1000000 == 0 {
                    if crawler.NumPages == last_pagecount {
                        noIncrease +=1
                    }
                    last_pagecount = crawler.NumPages
                }

                if noIncrease > 7 || time.Since(t0) >= crawler.MAX_TIME || crawler.NumPages >= crawler.MAX_PAGES {

                    glog.Info("Terminating crawler on a specified condition (time/no URLs left to crawl/ reached max)")
                    glog.Info("Total time spent crawling is ", time.Since(t0))
                    fmt.Printf("Status Update. Pages collected %d. Visited %d.\n", crawler.NumPages, len(vList))
                    fmt.Println("Total time: ", time.Since(t0))

                    // Tell workers to quit
                    for i:= 0; i< crawler.NumWorkers; i++ {
                        shutdown <- true
                    }

                    // Wait for workers to quit
                    wg.Wait()

                    // Close all channels
                    close(rurls)
                    close(surls)
                    close(shutdown)
                    close(pages)
                    fmt.Println("Done\n\n")
                    return nil
                }
        }
    }
}


// Prints out the sitemap of the crawler
func ( crawler *SingleCrawler ) Print() error {

    if err1 := IsOk(crawler); err1 != nil{
        return err1
    }

    stdout := os.Stdout
    outfile := stdout
    duped := true

    outfile, err := os.OpenFile( crawler.Filename, os.O_WRONLY | os.O_CREATE, 0644 )
    if err != nil {
        glog.Error("Unable to open requested file for writing. Defaulting to std out.")
        duped = false
    } else{
        os.Stdout = outfile
    }

    fmt.Printf("SiteMap from starting URL %s, total pages found %d.\n\n\n", crawler.Site.String(), crawler.NumPages )
    for i := 0; i < crawler.NumPages; i++ {
        crawler.Sitemap[i].Print(crawler.PRINT_LIMIT)
    }

    if duped == true {
        outfile.Close()
        os.Stdout = stdout
    }

    return nil

}


// Runs the crawler
func (mycrawler *SingleCrawler) Run() (error) {

    if err := IsOk(mycrawler); err != nil{
        return err
    }

    start := time.Now()

    // Start the crawler
    glog.Info("Starting web crawler")
    fmt.Println("starting web crawler")
    mycrawler.Start()

    // When done, print out the site map
    glog.Info("Done crawling, printing Sitemap")
    err := mycrawler.Print()
    if err!= nil{
        return errors.New("Unable to print sitemap")
    }

    // Log info when done, including elapsed time
    elapsed := time.Since(start)
    glog.Info( fmt.Sprintf("Finished Crawling Site, total elapsed time is %s", elapsed))
    fmt.Sprintf("Finished Crawling Site, total elapsed time is %s\n\n", elapsed)
    return nil

}
