package fsagent

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var size int64 = 200 * 1024 * 1024

func serveHTTP(conf Config, fnChannel chan<- string) {
	http.HandleFunc(conf.Folder, func(w http.ResponseWriter, r *http.Request) {
		var path string
		if err := r.ParseMultipartForm(size); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				file, _ := fileHeader.Open()
				path = fmt.Sprintf("%s", fileHeader.Filename)
				buf, _ := ioutil.ReadAll(file)
				tempFile, err := ioutil.TempFile("", fileHeader.Filename)
				if err != nil {
					fmt.Println(err)
				}
				tempFile.Write(buf)
				fnChannel <- tempFile.Name()
			}
		}
		fmt.Printf("File \"%v\" uploaded\n", path)
	})
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	fmt.Print(http.ListenAndServe(conf.Port, nil))
}
