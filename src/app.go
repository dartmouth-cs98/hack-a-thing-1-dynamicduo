package main

// https://github.com/disintegration/imaging
import "github.com/disintegration/imaging"

import (
  "fmt"
  "image"
  "image/gif"
  "os"
  "bytes"
  "log"
  "net/http"
  "io"
  "time"
  "crypto/md5"
  "strconv"
  "html/template"
  _ "image/jpeg"
  _ "image/png"
)

func main() {
  //use FileServer and Handle function to set up port to Listening
  // code taken from tutorial at http://www.alexedwards.net/blog/serving-static-sites-with-go
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    http.HandleFunc("/output", generateGIF)

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}

// upload logic
// adapted from 'upload' function at https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html
func generateGIF(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
     crutime := time.Now().Unix()
     h := md5.New()
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

     f1out, err := os.Create(handler1.Filename)
     if err != nil {
         fmt.Println(err)
         return
     }
     io.Copy(f1out, file1)
     if err != nil {
       fmt.Println("err: ", err)
     }

     // adapted from: https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang
     // with help from https://www.devdungeon.com/content/working-images-go
     f1p, err := os.Open(handler1.Filename)
     f1data, f1type, err := image.Decode(f1p)
     if err != nil {
       fmt.Println("err: %v", err)
     }
     fmt.Println(f1type)


     f2out, err := os.Create(handler2.Filename)
     if err != nil {
         fmt.Println(err)
         return
     }
     io.Copy(f2out, file2)

     f2p, err := os.Open(handler2.Filename)
     f2data, f2type, err := image.Decode(f2p)
     if err != nil {
       fmt.Println("err: %v", err)
     }
     fmt.Println(f2type)


     f3out, err := os.Create(handler3.Filename)
     if err != nil {
         fmt.Println(err)
         return
     }
     io.Copy(f3out, file3)

     f3p, err := os.Open(handler3.Filename)
     f3data, f3type, err := image.Decode(f3p)
     if err != nil {
       fmt.Println("err: ", err)
     }
     fmt.Println(f3type)

     //from tutorial at http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/
     // load static image and construct outGif
     outGif := &gif.GIF{}
     files := []image.Image {f1data, f2data, f3data}

     // figure out largest width and largest height of images uploaded by user
     // so that we can scale the images all to same width & height
     // otherwise, we cannot merge them into 1 GIF
     maxWidth := 0.0
     maxHeight := 0.0
     for _, f := range files {
       dx := float64(f.Bounds().Dx())
       dy := float64(f.Bounds().Dy())
       if dx > maxWidth {
         maxWidth = dx
       }
       if dy > maxHeight {
         maxHeight = dy
       }
     }

     //https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go for encoding other file formats like jpg/png
     //http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/ adapted using tutorial above
     //changed os to imaging
     for _, f := range files {

         buf := bytes.Buffer{}

         widthRatio := maxWidth/float64(f.Bounds().Dx())
         heightRatio := maxHeight/float64(f.Bounds().Dy())

         fscaled := ScaleImage(widthRatio, heightRatio, f, false)
         // fmt.Println("width: ", fscaled.Bounds().Dx())
         // fmt.Println("height: ", fscaled.Bounds().Dy())
         gif.Encode(&buf, fscaled, nil)
         gifGif, _ := gif.Decode(&buf)

         outGif.Image = append(outGif.Image, gifGif.(*image.Paletted))
         outGif.Delay = append(outGif.Delay, 0)
     }

     // save to out.gif
     //https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go
     final, err := os.Create("out.gif")
     if err != nil {
     		log.Fatalf("Erroradkslfjalsdfjlaksdf")
   	 }
     err = gif.EncodeAll(final, outGif)
     if err != nil {
       log.Fatalf("Erroradkslfjalsdfjlaksdf", err)
     }

     writeImage(w, outGif)
     final.Close()
  }
}

// writeImage adapted from https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang:
// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *gif.GIF) {

	buffer := new(bytes.Buffer)
	if err := gif.EncodeAll(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}


// following function adapted from: https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go
func ScaleImage(widthScale float64, heightScale float64, img image.Image, verbose bool) image.Image {

	newwidth := int(float64(img.Bounds().Dx()) * widthScale)
	newheight := int(float64(img.Bounds().Dy()) * heightScale)

	if verbose {
		log.Printf("Scaling image from (%d, %d) -> (%d, %d)", img.Bounds().Dx(), img.Bounds().Dy(), newwidth, newheight)
	}
	img = imaging.Resize(img, newwidth, newheight, imaging.Lanczos)

	return img
}
