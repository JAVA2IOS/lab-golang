package uploadController

import (
	"net/http"
	"log"
	"lab-golang/config"
	"os"
	"io"
	"mime/multipart"
	"errors"
)

type LabFileType string

const (
	LabFileTypeUnknown LabFileType = ""
	LabFileTypeTxt LabFileType = "txt"
	LabFileTypeXlsx LabFileType = "xlsx"
	LabFileTypeCSV LabFileType = "csv"
)

func SaveFile(w http.ResponseWriter, r *http.Request, fileType LabFileType) string {
	err := make(chan error)
	config.HandlerHttpError(w, err)

	r.ParseForm()
	//把上传的文件存储在内存和临时文件中
	err <- r.ParseMultipartForm(32<<20)
	xlsxFile, xlsxHandler, xlsxErr := r.FormFile(string(fileType))
	err <- xlsxErr
	defer xlsxFile.Close()

	filePath, fileErr := saveFileToLocalPath(xlsxHandler.Filename, xlsxFile, xlsxHandler)
	err <- fileErr

	return filePath
}

func saveFileToLocalPath(fileName string, file multipart.File, header *multipart.FileHeader) (string, error) {
    //创建上传的目的文件
    filePath := "./upload/" + fileName

    f, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE, 0666)

    if err != nil {
        log.Printf("file open failed : %v \n", err.Error())
        return "", errors.New("文件创建失败: " + err.Error())
    }
    defer f.Close()

    //拷贝文件
    _, err = io.Copy(f, file)
    if err != nil {
    	log.Printf("file copied failed : %v \n", err.Error())
    	return "", errors.New("文件数据上传失败: %v" + err.Error())
    }

    return filePath, nil
}
