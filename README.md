# petitcrawler
Info: single domain web crawler package, implemented in Golang
Author: Rachel Mendiola, June 6 2017

Description:
petitcrawler crawls a single domain and outputs a sitemap to a file. 
The basic input is a starting URL to crawl, which needs to look something like "http://notrealURL.com".
Written as a technical challenge for a company. My first Go program!
If you have any suggestions or comments please make them known to me.


Dependent on (install with 'go get'):
- "golang.org/x/net/html"
- "github.com/golang/glog"
- "github.com/asaskevich/govalidator"


The main idea:
- Two main entites: Crawler and Worker
- Crawler spawns a set number of goroutine workers
- Crawler sends urls to the workers
- Workers fetch urls, parse html, create Pages
- Workers pass back to the Crawler new urls found that are within the original domain
- Workers pass back to the Crawler new Pages that are created
- Crawler accepts the new urls, makes sure they are unique
- Crawler accepts the new Pages, makes sure they are unique
- Crawler watches over several termination conditions and sends 'shutdown' signal to workers
There are some UML documents in doc/ if you want to see this description drawn up. 

This package is contained in petitcrawler/. 
A test suite is included in petitcrawlertest/.
An example program using the package is in test/.


There are many command-line options, to show those: ./test -help
A command line example is: ./test -url <URL> -numworkers=100 -maxtime=60
