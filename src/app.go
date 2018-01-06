package main

import (
  "log"
  "net/http"
)

func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}

// use FileServer and Handle function to set up port to Listening
// code taken from tutorial at http://www.alexedwards.net/blog/serving-static-sites-with-go
