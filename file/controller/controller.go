package controller

import (
	"net/http"
	"os"
	"log"
	"io/ioutil"
)

func WebOutputFile(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Length", value)
	
}

func ClearFile(filePath string) error {
	dir_list, e := ioutil.ReadDir(filePath)
    if e != nil {
    	log.Printf("文件路径错误: %v", filePath)
        return e
    }

    for _, v := range dir_list {
    	err := os.Remove(filePath + "/" + v.Name())
    	if err != nil {
    		log.Printf("删除文件出错: %v", err.Error())
    		return err
    		break
    	}else {
    		log.Printf("删除文件[%v]成功", v.Name())
    	}
    }

    return nil
}