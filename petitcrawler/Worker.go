package petitcrawler


import (
    "fmt"
    "net/http"
    "net/url"
    "golang.org/x/net/html"
    "time"
    "strings"
    "sync"
    "errors"
    "github.com/golang/glog"
    "github.com/asaskevich/govalidator"
)



// One Worker process. Accepts urls in channel url. Accepts termination signal in shutdown.
// Process url received, send back to controller in send_back.
// Send back successfully crawled page data to controller. 
func Worker( myID int, urls chan string, send_back chan string, domain *url.URL, pages chan Page, shutdown <- chan bool, wg *sync.WaitGroup ) {

    defer wg.Done()
    defer glog.Flush()

    for {
        select {
            case _ = <- shutdown:
                return
            case link := <- urls:
                p, err := Work( link, send_back, domain )
                if err == nil {
                    select{
                        case <-time.After(5*time.Second):
                        case pages <- p:
                    }
                } 
        }
    }

}



// Worker makes an http Get request to the given URL and parses the body of the html doc
// using a separate recursive function
// @Return is a create Page (urls, assets) and an integer 0 for success, -1 for fail
func Work( link string, uList chan string, domain *url.URL ) (Page, error) {

    t0 := time.Now()
    var page Page

    if govalidator.IsURL(link) == false {
        return page, errors.New( fmt.Sprintf("Not a url %s.",link))
    }
    if err:= DomainCheck(domain); err!= nil{
        return page, err
    }
    
    // Make a request 
    glog.Info( fmt.Sprintf("Requesting to URL %s.", link ) )
    //timeout := time.Duration( 6 * time.Second)
    //client := http.Client{ Timeout: timeout, } 
    resp, err := http.Get(link)
    
    if err != nil {

        // Try one more time, but be respectful of websites! Do not send too many requests.
        resp, err = http.Get(link)
        if err != nil {
            glog.Warning( fmt.Sprintf("No response from %s. Error is %s. Skipping URL.\n", link, err))
            return page, errors.New( fmt.Sprintf("No response form %s. Error is %s.", link, err))
        }
    }
    defer resp.Body.Close()
    
    // If we have an error, log it and return
    if resp.StatusCode != 200 {
        glog.Warning( fmt.Sprintf("Bad response from http request to page %s. Error code %d. Skipping URL\n", link, resp.StatusCode))
        return page, errors.New( fmt.Sprintf("Bad response code from request to page %s. Error code %d.", link, resp.StatusCode))
    }

    // Parse the body of the response
    doc, err := html.Parse( resp.Body )
    if err != nil {
        glog.Warning( fmt.Sprintf("Unable to parse html from page %s. Error is %s. Skipping URL.\n", link, err))
        return page, errors.New( fmt.Sprintf("Unable to parse html of page %s.", link))
    }

    // Search the html structure for links, static assets
    err = CheckNode( doc, uList, domain, &page, t0 )
    page.MyUrl = link

    glog.Info( fmt.Sprintf("Done crawling link %s\n", link))
    if err != nil{ 
        return page, errors.New( fmt.Sprintf("Error parsing html page %s.", link))
    } else {
        return page, nil
    }
}


// CheckNode searches one node in a parsed HTML tree, looking for 
// URLS and static assets to record. 
// Information found is passed back through the @param page *Page.
func CheckNode( n *html.Node, uList chan string, domain *url.URL, page *Page, t0 time.Time) error {
    
    if n == nil {
        return nil 
    }
    if page == nil{
        return errors.New("Bad page struct pointer.")
    }
    if err:= DomainCheck(domain); err!= nil{
        return err
    }

    // Search for links, images, scripts
    if n.Type == html.ElementNode && ( n.Data == "a" || n.Data == "img" || n.Data == "link" || n.Data == "script") { 


        // Iterate through the attributes of each html element node
        for _, a := range n.Attr { 
            if time.Since(t0) >= time.Duration(5)*time.Second{ return errors.New("Timeout")}

            // Record sources of images, links, and scripts
            if a.Key == "src" {
                page.Assets = append( page.Assets, a.Val )
                break
            }

            // Record new URLs found, and send them back to the parent crawler
            if a.Key == "href"  { 
                // If we can't parse the URL and find the domain, skip this URL, since we can't understand it
                str := a.Val
                if len(a.Val) >0 && a.Val[0] == '.' {
                    str = str[1:]
                }
                u,err := url.Parse( str )
                if err != nil {
                    continue
                }

                if u.Host == "" {
                    u.Host = domain.Host
                }
                if u.Scheme == "javascript" {
                    page.Assets = append( page.Assets, a.Val )
                    break
                }
                if strings.Contains(u.String(), ".png") || strings.Contains(u.String(), ".jpg")|| strings.Contains(u.String(), ".ico") || strings.Contains(u.String(), ".css") {
                    page.Assets = append( page.Assets, a.Val )
                    break
                }

                // Check to see if the discovered URL is within the original domain
                if (u.Host == domain.Host || u.Host == domain.Path) && strings.Contains(u.String(), "mailto") == false {
                    url := a.Val

                    // Check to see if the URL is absolute, if not, specify a scheme
                    if u.IsAbs() == false {
                        u.Scheme = domain.Scheme
                        url = u.String()
                    }

                    // Record this info in the Page, send to crawler with URL channel
                    page.BabyUrls = append( page.BabyUrls, url )
                    select{ 
                        case <-time.After(2*time.Second):
                            return errors.New("Timeout waiting for write to channel") 
                        case uList <- url:
                    }
                }
                break
            }
        }
    }

    // Recursively iterate over all nodes in the html parse tree
    for c := n.FirstChild; c != nil; c = c.NextSibling { 
        err := CheckNode(c, uList, domain, page, t0)
        if err != nil { return err}
    }
    return nil
}
