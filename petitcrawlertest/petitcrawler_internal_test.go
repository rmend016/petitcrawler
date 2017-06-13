package petitcrawler_test


import (
    "petitcrawler"
    "testing"
    "time"
    "os"
    "net/url"
    "golang.org/x/net/html"
    "sync"
    "fmt"
)


// Unit test DomainCheck with many test cases
func TestDomainCheck( t *testing.T) {
    d, err := url.Parse("http://google.com")
    if err!= nil{ return }
    if err = petitcrawler.DomainCheck(d); err != nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 1.")
    }
    d, _ = url.Parse("")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 2.")
    }
    d, _ = url.Parse("https://")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 3.")
    }
    d, _ = url.Parse("https://yellow")
    if err = petitcrawler.DomainCheck(d); err != nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 4.")
    }
    d, _ = url.Parse("noscheme.com")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 5.")
    }
    d, _ = url.Parse("www.google.com")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 6.")
    }
    d, _ = url.Parse("http://www.google.com")
    if err = petitcrawler.DomainCheck(d); err != nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 7.")
    }
    d, _ = url.Parse(".com")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 8.")
    }
    d, _ = url.Parse("𝔣𝔡𝔰𝔣")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 9.")
    }
    d, _ = url.Parse("\n\n")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 10.")
    }
    d, _ = url.Parse("ᵃ𝐋ιⒸ𝒾Ø𝓤𝐬")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 11.")
    }
    d, _ = url.Parse(" ｷo尺 ﾘou尺 g尺ouｱ ｲﾑg")
    if err = petitcrawler.DomainCheck(d); err == nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 12.")
    }
    d, _ = url.Parse("http://ko.wikipedia.org/wiki/%EC%9C%84%ED%82%A4%EB%B0%B1%EA%B3%BC:%EB%8C%80%EB%AC%B8")
    if err = petitcrawler.DomainCheck(d); err != nil{
        t.Fatalf("TestDomainCheck() Failed: Test case 13.")
    }
    
}


// Unit test CheckNode bad domain
func TestCheckNodeBadDomain( t *testing.T) {
    code, err := os.Open("htmlexample.txt")
    doc, err := html.Parse(code)
    var page petitcrawler.Page
    t0 := time.Now()
    domain, _ := url.Parse("http://")
    uList := make( chan string )
    err = petitcrawler.CheckNode( doc, uList, domain, &page, t0)
    if err == nil {
        t.Fatalf("TestCheckNodeBadDomain() Failed: Expecting to fail on bad domain.")
    }
}


// Unit test CheckNode nil page *Page
func TestCheckNodeBadPageStruct( t *testing.T) {
    code, err := os.Open("htmlexample.txt")
    doc, err := html.Parse(code)
    var page *petitcrawler.Page
    t0 := time.Now()
    domain, _ := url.Parse("http://google.com")
    uList := make( chan string )
    err = petitcrawler.CheckNode( doc, uList, domain, page, t0)
    if err == nil {
        t.Fatalf("TestCheckNodeBadPageStruct() failed: Expecting to fail on bad page ptr.")
    }
}


// Unit test CheckNode nil html node
func TestCheckNodeBadHtmlNode( t *testing.T) {
    var doc *html.Node
    var page petitcrawler.Page
    t0 := time.Now()
    domain, _ := url.Parse("http://google.com")
    uList := make( chan string, 100 )
    err := petitcrawler.CheckNode( doc, uList, domain, &page, t0)
    if err != nil {
        t.Fatalf("TestCheckNodeBadHtmlNode() failed: %s. Expecting to succeed.", err)
    }
}


// Unit test Work bad domain
func TestWorkBadDomain( t *testing.T) {
    domain, _ := url.Parse("http://")
    uList := make( chan string, 100 )
    _, err :=  petitcrawler.Work( "http://google.com/search", uList, domain)
    if err == nil {
        t.Fatalf("TestWorkBadDomain() failed: %s. Expecting to succeed.", err)
    }
}


// Unit test Work basic parse 
func TestWorkGoodChan( t *testing.T) {
    domain, _ := url.Parse("http://google.com")
    uList := make( chan string, 100 )
    _, err :=  petitcrawler.Work( "http://google.com/search", uList, domain)
    if err != nil {
        t.Fatalf("TestWorkGoodChan() failed: %s. Expecting to succeed.", err)
    }
}


// Unit test Work returns on full/block channel
func TestWorkBadChan(t *testing.T) {
    domain, _ := url.Parse("http://google.com")
    uList := make( chan string )
    _, err :=  petitcrawler.Work( "http://google.com/search", uList, domain)
    if err == nil {
        t.Fatalf("TestWorkBadChan() failed: %s. Expecting to fail on full channel.", err)
    }
}


// Unit test Work bad link
func TestWorkBadLink(t *testing.T) {
    domain, _ := url.Parse("http://google.com")
    uList := make( chan string )
    _, err :=  petitcrawler.Work( "brokenlink", uList, domain)
    if err == nil {
        t.Fatalf("TestWork() failed: %s. Expecting to fail on broken url.", err)
    }
}


// Test completion on a full, blocked channel, and empty
func TestWork(t *testing.T) {
    domain, _ := url.Parse("http://google.com")
    
    uList := make( chan string )
    _, err :=  petitcrawler.Work( "http://google.com/search", uList, domain)
    if err == nil {
        t.Fatalf("TestWork() failed: %s. Expecting to block on channel and exit", err)
    }

    uList2 := make( chan string, 100)
    _, err =  petitcrawler.Work( "http://google.com/search", uList2, domain)
    if err != nil {
        t.Fatalf("TestWork() failed: %s. Expecting to have enough space", err)
    }
}


// Unit test Worker exits normally given work
func TestWorkerShutdown(t *testing.T) {
    numw := petitcrawler.MAX_WORKERS
    urls := make( chan string, numw)
    rurls := make( chan string, numw)
    pages := make( chan petitcrawler.Page, numw)
    shutdown := make( chan bool, numw)
    domain, _ := url.Parse("http://hi.com")
    var wg sync.WaitGroup
    for i:=0; i<numw; i++ {
        wg.Add(1)
        go petitcrawler.Worker( i , urls, rurls, domain, pages, shutdown, &wg )
    }
    urls <- "http://hi.com"
    urls <- "http://google.com"
    urls <- "http://hello.com"
    for i:=0; i<numw; i++ {
        shutdown <- true
    }
    wg.Wait()
}


// Unit test Worker exits normally given no work
func TestWorkerShutdownNoWork(t *testing.T) {
    numw := petitcrawler.MAX_WORKERS
    urls := make( chan string, numw)
    rurls := make( chan string, numw)
    pages := make( chan petitcrawler.Page, numw)
    shutdown := make( chan bool, numw)
    domain, _ := url.Parse("http://hi.com")
    var wg sync.WaitGroup
    for i:=0; i<numw; i++ {
        wg.Add(1)
        go petitcrawler.Worker( i , urls, rurls, domain, pages, shutdown, &wg )
    }
    for i:=0; i<numw; i++ {
        shutdown <- true
    }
    wg.Wait()
}


// Unit test Worker exits on blocked/full channel pages
func TestWorkerShutdownBadChanPages(t *testing.T) {
    numw := petitcrawler.MAX_WORKERS
    urls := make( chan string, numw)
    rurls := make( chan string, numw)
    pages := make( chan petitcrawler.Page)
    shutdown := make( chan bool, numw)
    domain, _ := url.Parse( "http://hi.com" )
    var wg sync.WaitGroup
    for i:=0; i<numw; i++ {
        wg.Add(1)
        go petitcrawler.Worker( i , urls, rurls, domain, pages, shutdown, &wg )
    }
    urls <- "http://hi.com"
    urls <- "http://google.com"
    urls <- "http://hello.com"
    for i:=0; i<numw; i++ {
        shutdown <- true
    }
    wg.Wait()
}


// Unit test Worker exits on blocked/full channel rurls
func TestWorkerShutdownBadChanRurls(t *testing.T) {
    numw := petitcrawler.MAX_WORKERS
    urls := make( chan string, numw)
    rurls := make( chan string)
    pages := make( chan petitcrawler.Page, numw)
    shutdown := make( chan bool, numw)
    domain, _ := url.Parse( "http://hi.com" )
    var wg sync.WaitGroup
    for i:=0; i<numw; i++ {
        wg.Add(1)
        go petitcrawler.Worker( i , urls, rurls, domain, pages, shutdown, &wg )
    }
    urls <- "http://hi.com"
    urls <- "http://google.com"
    urls <- "http://hello.com"
    for i:=0; i<numw; i++ {
        shutdown <- true
    }
    wg.Wait()
}


// Unit test Worker exits, all channels blocked
func TestWorkerShutdownBadChanAll(t *testing.T) {
    numw := petitcrawler.MAX_WORKERS
    urls := make( chan string, numw)
    rurls := make( chan string)
    pages := make( chan petitcrawler.Page)
    shutdown := make( chan bool, numw)
    domain, _ := url.Parse("http://hi.com")
    var wg sync.WaitGroup
    for i:=0; i<numw; i++ {
        wg.Add(1)
        go petitcrawler.Worker( i , urls, rurls, domain, pages, shutdown, &wg )
    }
    for i:=0; i<numw; i++ {
        shutdown <- true
    }
    wg.Wait()
}


// Unit test Page.Print @limit printing
func TestPagePrint(t *testing.T){
    var p petitcrawler.Page
    p.MyUrl = "http://google.com"
    p.Assets = append(p.Assets, "first asset")
    p.BabyUrls = append( p.BabyUrls, "url found on this page")
    for i:= 0; i < 1000; i++ {
        p.Assets = append(p.Assets, fmt.Sprintf("asset#%d", i))
        p.BabyUrls = append( p.BabyUrls, fmt.Sprintf("url#%d",i))
    }
}


// Unit test Run with bad crawler Site
func TestRunBadSite(t *testing.T){
    c, err := petitcrawler.New()
    if err != nil {
        t.Fatalf( fmt.Sprintf("TestRunBadSite() Failed to create crawler. %s.", err))
    }
    c.Site = nil
    err = petitcrawler.Run(c)
    if err == nil{
        t.Fatalf("TestRunBadSite() Failed: Expecting to quit on crawler with no Site.")
    }
}


// Unit test Run with bad crawler Sitemap
func TestRunBadSitemap(t *testing.T){
    c, err := petitcrawler.New()
    if err != nil {
        t.Fatalf( fmt.Sprintf("TestRunBadSitemap() Failed to create crawler. %s.", err))
    }
    c.Sitemap = nil
    err = petitcrawler.Run(c)
    if err == nil{
        t.Fatalf("TestRunBadSitemap() Failed: Expecting to quit on crawler with no Sitemap.")
    }
}


// Unit test Run with bad crawler num workers <= 0
func TestRunNegativeWorkers(t *testing.T){
    c, err := petitcrawler.New()
    if err != nil {
        t.Fatalf( fmt.Sprintf("TestRunBadNegativeWorkers() Failed to create crawler. %s.", err))
    }
    c.NumWorkers = -2000
    err = petitcrawler.Run(c)
    if err == nil{
        t.Fatalf("TestRunNegativeWorkers() Failed: Expecting to quit on crawler with negative # workers.")
    }
}


// Unit test Run with bad crawler no filename
func TestRunBadFilename(t *testing.T){
    c, err := petitcrawler.New()
    if err != nil {
        t.Fatalf( fmt.Sprintf("TestRunBadFilename() Failed to create crawler. %s.", err))
    }
    c.Filename = ""
    err = petitcrawler.Run(c)
    if err == nil{
        t.Fatalf("TestRunBadFilename() Failed: Expecting to quit on crawler with no filename.")
    }
}


// Unit test SingleCrawler.Print
func TestPrint(t *testing.T) {
    c, err := petitcrawler.New()
    if err != nil {
        t.Fatalf( fmt.Sprintf("TestPrint() Failed to create crawler. %s.", err))
    }
    err = c.Print()
    if err != nil {
        t.Fatalf( fmt.Sprintf("TestPrint() Failed to print crawler sitemap. %s.", err))
    }
}
