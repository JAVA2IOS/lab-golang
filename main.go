package main

import (
	"lab-golang/router"
)


func main() {

	// index
	// http.HandleFunc("/", HTMLPageRouter)
	// http.HandleFunc("/index", HTMLPageRouter)

	// xlsx
	// http.HandleFunc("/api/v0/file/upload/xlsx", file.UploadXlsxFile)

	// NetResponseHandler()
	netRouter.NetResponseHandler()
}
