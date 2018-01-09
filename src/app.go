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
  "encoding/base64"
  _ "image/jpeg"
  _ "image/png"
)

// http.HandleFunc("/upload", upload)

func main() {
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

     // fmt.Fprintf(w, "%v", handler1.Header)

     // f1, err := os.OpenFile(handler1.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     f1out, err := os.Create(handler1.Filename)
     // f1, err := os.Open(handler1.Filename)
     if err != nil {
         fmt.Println(err)
         return
     }
     // defer f1.Close()
     io.Copy(f1out, file1)
     if err != nil {
       fmt.Println("err: ", err)
     }

     // adapted from: https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang
     // m := image.NewRGBA(image.Rect(0, 0, 240, 240))
     // with help from https://www.devdungeon.com/content/working-images-go
     // f1out.seek(0,0)
     f1p, err := os.Open(handler1.Filename)
     f1data, f1type, err := image.Decode(f1p)
     if err != nil {
       fmt.Println("err: %v", err)
     }

     fmt.Println(f1type)
     // writeImage(w, &f1data)
     // fmt.Fprintf(w, "%v", handler2.Header)
     // createf2, err := os.OpenFile(handler2.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     f2out, err := os.Create(handler2.Filename)
     // f2, err := os.Open(handler2.Filename)
     if err != nil {
         fmt.Println(err)
         return
     }
     // defer f2.Close()
     io.Copy(f2out, file2)

     // f2out.seek(0,0)
     f2p, err := os.Open(handler2.Filename)
     f2data, f2type, err := image.Decode(f2p)
     if err != nil {
       fmt.Println("err: %v", err)
     }
     fmt.Println(f2type)

     // fmt.Fprintf(w, "%v", handler3.Header)
     // createf3, err := os.OpenFile(handler3.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     f3out, err := os.Create(handler3.Filename)
     // f3, err := os.Open(handler3.Filename)
     if err != nil {
         fmt.Println(err)
         return
     }
     // defer f3.Close()
     io.Copy(f3out, file3)
     // f3out.seek(0,0)
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

     //https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go for encoding other file formats like jpg/png
     //http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/ adapted using tutorial above
     //changed os to imaging
     totalWidth := 0
     totalHeight := 0
     for _, f := range files {
       totalWidth = totalWidth + f.Bounds().Dx()
       totalHeight = totalHeight + f.Bounds().Dy()
     }

     for _, f := range files {
         // f, _ := imaging.Open(name)
         // inGif := &gif.GIF{}
         buf := bytes.Buffer{}

         // fscaled := ScaleImage(10, f, false)
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
     err = gif.EncodeAll(final, outGif)
     if err != nil {
       log.Fatalf("Erroradkslfjalsdfjlaksdf", err)
     }

     writeImage(w, outGif)
     final.Close()
  }
}

// from https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang
var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/gif;base64,{{.Image}}"></body>`

// from https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang:
// Writeimagewithtemplate encodes an image 'img' in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := gif.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}

// from https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang:
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


// https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go
func ScaleImage(scale float64, img image.Image, verbose bool) image.Image {
	//Scale operation. Ignore if scale is 1.0
	if scale != 1.0 {
		newwidth := int(float64(img.Bounds().Dx()) * scale)
		newheight := int(float64(img.Bounds().Dy()) * scale)

		if verbose {
			log.Printf("Scaling image from (%d, %d) -> (%d, %d)", img.Bounds().Dx(), img.Bounds().Dy(), newwidth, newheight)
		}
		img = imaging.Resize(img, newwidth, newheight, imaging.Lanczos)
	}
	return img

}
