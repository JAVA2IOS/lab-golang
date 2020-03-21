package netRouter

import (
	"html/template"
	"log"
	"net/http"
	qql "lab-golang/excel/qql/controller"
	"lab-golang/tool/yaml"
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
		qql.QqlCreateShippingOrder(w, r)
	case "/qql/tool/excel/clearfiles":
		qql.QqlClearXlsxFiles(w, r)
	case "/qql/tool/excel/xlsxlist":
		qql.QqlXlsxFileList(w, r)
	default:
		log.Printf("当前路径失败: %v", r.URL.Path)
		tpl, gloableErr := template.ParseFiles("html/404.html")
		if gloableErr != nil {
			http.NotFound(w, r)
			return
		}
		tpl.Execute(w, nil)
	}
}

func NetResponseHandler() {
	reader, configureErr := yamlReader.Instance()
	if configureErr != nil {
		log.Fatal("yaml文件配置失败: ", configureErr, "code: ", configureErr.Error())
	}

	downloadPath := "/" + reader.Configure.Xlsx.DownloadFile + "/"
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	http.Handle(downloadPath, http.StripPrefix(downloadPath, http.FileServer(http.Dir("file/xlsx"))))

	http.HandleFunc("/", netRouterConfigureCenter)

	err := http.ListenAndServe(":" + reader.Configure.Port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err, "code: ", err.Error())
	}
}
