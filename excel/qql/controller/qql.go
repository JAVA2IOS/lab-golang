package controller

import (
	"net/http"
	"lab-golang/upload/controller"
	"lab-golang/config"
	"path/filepath"
	"log"
	file "lab-golang/file/controller"
	"lab-golang/excel/qql/model"
)

func QqlCreateShippingOrder(w http.ResponseWriter, r *http.Request) {
	filePath, fileName, err := uploadController.SaveFile(w, r, "csv_file")
	if err != nil {
		config.HandlerHttpError(w, err)
		return
	}

	data , err := uploadController.ReadCSVFile(filePath, fileName)
	if err != nil {
		config.HandlerHttpError(w, err)
		return
	}


	newFilePath, err := uploadController.CreateNewShippingOrder(data)
	if err != nil {
		config.HandlerHttpError(w, err)
		return 
	}

	config.HandlerHttpJson(w, "../" + newFilePath)
}

func QqlClearXlsxFiles(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs("file/xlsx")
	if err != nil {
		config.HandlerHttpError(w, err)
		return 
	}
	log.Printf("地址: %v", path)

	// 清除xlsx文件
	err = file.ClearFile(path)
	if err != nil {
		config.HandlerHttpError(w, err)
		return 
	}

	sourcePath, sourceErr := filepath.Abs("upload/file")
	if sourceErr != nil {
		config.HandlerHttpError(w, sourceErr)
		return
	}
	log.Printf("地址: %v", sourcePath)

	// 清除xlsx文件
	sourceErr = file.ClearFile(sourcePath)
	if sourceErr != nil {
		config.HandlerHttpError(w, sourceErr)
		return 
	}

	config.HandlerHttpJson(w, "清除文件成功")
}

func QqlXlsxFileList(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs("file/xlsx")
	if err != nil {
		config.HandlerHttpError(w, err)
		return 
	}
	log.Printf("地址: %v", path)

	lists := file.ReadFileList(path)

	if len(lists) == 0 {
		config.HandlerHttpJson(w, nil)
		return 
	}

	xlsxFile := make([]*model.XlsxFile, 0)
	for _, v := range lists {
		file := &model.XlsxFile{Name: v.Name(), Size: v.Size(), CreateTime: v.ModTime().Unix()}
		xlsxFile = append(xlsxFile, file)
	}

	config.HandlerHttpJson(w, xlsxFile)
}








