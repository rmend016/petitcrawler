package petitcrawler 


import (
    "fmt"
)


// Holds information about pages crawled,
// keeps track of URL, Assets, and any other URLs found on the page.
type Page struct { 

    MyUrl string        // the URL of the Page
    Assets []string     // static Assets
    BabyUrls []string    // the URL of the Page this link was found on

}


// Prints out the information a Page struct is storing
// Only prints the first PRINT_LIMIT Assets and URLS
func ( page *Page ) Print(PRINT_LIMIT int) {

    fmt.Printf( "Page URL: %s\n\n", page.MyUrl )
    if len( page.Assets ) > PRINT_LIMIT {
        fmt.Printf( "Assets (%d):\n\t%s\n\n", len(page.Assets), page.Assets[0:PRINT_LIMIT] )
    } else { 
        fmt.Printf( "Assets (%d):\n\t%s\n\n", len(page.Assets), page.Assets )
    }
    if len( page.BabyUrls ) > PRINT_LIMIT {
        fmt.Printf( "Children URLs (%d):\n\t%s\n\n\n\n", len(page.BabyUrls), page.BabyUrls[0:PRINT_LIMIT] )
    } else {
        fmt.Printf( "Children URLs (%d):\n\t%s\n\n\n\n", len(page.BabyUrls), page.BabyUrls )
    }

}
