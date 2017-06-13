package main

import (
    "petitcrawler"
    "os"
    "fmt"
)

func main() {

    Mycrawler, err := petitcrawler.New()
    if err != nil {
        fmt.Println("Failed to create crawler, error is: ", err)
        os.Exit(1)
    }

//    err = Mycrawler.Start()
    err = Mycrawler.Run()
    if err !=nil {
        fmt.Println("Failed to run crawler, error is: ", err)
        os.Exit(1)
    }
    //err = Mycrawler.Print()
}
