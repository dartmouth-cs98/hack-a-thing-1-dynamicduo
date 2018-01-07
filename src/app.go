package main

import "fmt"
import "image"
import "image/gif"
import "os"
import "bytes"
// https://github.com/disintegration/imaging
import "github.com/disintegration/imaging"

import (
  "log"
  "net/http"
  "io"
  "time"
  "crypto/md5"
  "strconv"
  "html/template"
)

// http.HandleFunc("/upload", upload)

func main() {

  //from tutorial at http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/
  files := []string{"g1.gif", "g2.gif","g3.gif"}

    // load static image and construct outGif
    outGif := &gif.GIF{}

    //https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go for encoding other file formats like jpg/png
    //http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/ adapted using tutorial above
    //changed os to imaging
    for _, name := range files {
        f, _ := imaging.Open(name)
        // inGif := &gif.GIF{}
        buf := bytes.Buffer{}

        gif.Encode(&buf, f, nil)
        gifGif, _ := gif.Decode(&buf)
        // f.Close()

        outGif.Image = append(outGif.Image, gifGif.(*image.Paletted))
        outGif.Delay = append(outGif.Delay, 0)
    }

    // save to out.gif
    //https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go
    final, err := os.Create("out.gif")
    if err != nil {
		log.Fatalf("Erroradkslfjalsdfjlaksdf")
	  }
    gif.EncodeAll(final, outGif)
    final.Close()

  //use FileServer and Handle function to set up port to Listening
  // code taken from tutorial at http://www.alexedwards.net/blog/serving-static-sites-with-go
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    http.HandleFunc("/upload", upload)

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)

}

// upload logic
// taken from https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html
func upload(w http.ResponseWriter, r *http.Request) {
  fmt.Println("fjdksalfjdsalf")
  fmt.Println("method:", r.Method)
  if r.Method == "GET" {
     crutime := time.Now().Unix()
     h := md5.New()

     fmt.Println("here")
     io.WriteString(h, strconv.FormatInt(crutime, 10))
     token := fmt.Sprintf("%x", h.Sum(nil))

     t, _ := template.ParseFiles("example.html")
     t.Execute(w, token)
  } else {
     r.ParseMultipartForm(32 << 20)
     file1, handler1, err1 := r.FormFile("file1")
     file2, handler2, err2 := r.FormFile("file2")
     file3, handler3, err3 := r.FormFile("file3")
     if err1 != nil {
         fmt.Println(err1)
         return
     }
     if err2 != nil {
         fmt.Println(err2)
         return
     }
     if err3 != nil {
         fmt.Println(err3)
         return
     }
     defer file1.Close()
     defer file2.Close()
     defer file3.Close()

     // fmt.Fprintf(w, "%v", handler1.Header)

     f1, err := os.OpenFile(handler1.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     if err != nil {
         fmt.Println(err)
         return
     }
     defer f1.Close()
     io.Copy(f1, file1)
     fmt.Fprintf(w, "<html><body><img src=%v /></body></html>", handler1.Header)

     // fmt.Fprintf(w, "%v", handler2.Header)
     f2, err := os.OpenFile(handler2.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     if err != nil {
         fmt.Println(err)
         return
     }
     defer f2.Close()
     io.Copy(f2, file2)

     // fmt.Fprintf(w, "%v", handler3.Header)
     f3, err := os.OpenFile(handler3.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     if err != nil {
         fmt.Println(err)
         return
     }
     defer f3.Close()
     io.Copy(f3, file3)
  }
}
