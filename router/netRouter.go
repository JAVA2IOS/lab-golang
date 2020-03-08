package netRouter

import (
	"html/template"
	"log"
	"net/http"
	// "fmt"
	qql "lab-golang/excel/qql/controller"
)

func netRouterConfigureCenter(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fallthrough
	case "/index":
		tpl, gloableErr := template.ParseFiles("html/lab.html")
		if gloableErr != nil {
			http.NotFound(w, r)
			return
		}
		tpl.Execute(w, nil)
	case "/qql/tool/excel/order":
		go qql.QqlCreateShippingOrder(w, r)
	default:
		tpl, gloableErr := template.ParseFiles("html/404.html")
		if gloableErr != nil {
			http.NotFound(w, r)
			return
		}
		tpl.Execute(w, nil)
	}
}

func NetResponseHandler() {
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))

	http.HandleFunc("/", netRouterConfigureCenter)

	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err, "code: ", err.Error())
	}
}
