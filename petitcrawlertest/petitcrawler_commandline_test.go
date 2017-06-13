package petitcrawler_test


import (
    "petitcrawler"
    "testing"
    "time"
)


// Unit test New creates crawler with command line flags
// must pass in -url command line arg
func TestNew(t *testing.T) {
    t.Parallel()
    Mycrawler, err := petitcrawler.New() 
    if err != nil {
        t.Fatalf("TestNew() failed: %s", err)
    }
    if Mycrawler.NumWorkers != *petitcrawler.NumwPtr {
        t.Fatalf("Incorrectly set num workers, should be %d, got: %d.", *petitcrawler.NumwPtr, Mycrawler.NumWorkers)
    }
    if Mycrawler.MAX_PAGES != *petitcrawler.MaxcPtr{
        t.Fatalf("Incorrectly set max pages, should be %d, got: %d.", *petitcrawler.MaxcPtr, Mycrawler.MAX_PAGES)
    }
    if Mycrawler.PRINT_LIMIT != *petitcrawler.MaxpPtr{
        t.Fatalf("Incorrectly set print limit, should be %d, got: %d.", *petitcrawler.MaxpPtr, Mycrawler.PRINT_LIMIT)
    }
    if Mycrawler.NumPages != 0 {
        t.Fatalf("Incorrectly set num pages, should be %d, got: %d.", 0, Mycrawler.NumPages)
    }
    if Mycrawler.MAX_TIME != (time.Duration(*petitcrawler.MaxtPtr))*(time.Second){
        t.Fatalf("Incorrectly set max time, should be %d, got: %d.", *petitcrawler.MaxtPtr, Mycrawler.MAX_TIME)
    }
    if len(Mycrawler.Sitemap) != *petitcrawler.MaxcPtr {
        t.Fatalf("Incorrectly initialized Sitemap of crawler, should be %d, got: %d.", *petitcrawler.MaxcPtr, len(Mycrawler.Sitemap))
    }
    if len(Mycrawler.Filename) >= 255 {
        t.Fatalf("Incorrectly initialized Filename of crawler, should be less that 255 in length, got: %d in length.", len(Mycrawler.Filename))
    }
}


// Unit test Run function
func TestRun(t *testing.T) {

    Mycrawler, err := petitcrawler.New()
    if err != nil {
        t.Fatalf("TestRun() failed: %s", err)
    }

    err = petitcrawler.Run(Mycrawler)
    if err !=nil {
        t.Fatalf("TestRun() failed. %s", err)
    }
}

