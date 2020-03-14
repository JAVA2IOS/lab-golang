package tool

import (
	"strings"
	"strconv"
	"path"
	"os"
)

func GetFileSuffixString(filePath string) string {
	file := path.Base(filePath)

	suffix := path.Ext(file)
	return suffix
}

func GetFileName(filePath string) string {
	file := path.Base(filePath)
	suffix := path.Ext(file)

	fileName := strings.TrimSuffix(file, suffix)

	return fileName
}

func GetFilePath(filePath string) string {
	file := path.Base(filePath)

	newPath := strings.TrimSuffix(filePath, file)

	return newPath
}

func AppendFileSuffix(filePath string, suffix string) string {
	return GetFileName(filePath) + "." + suffix
}


func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func NewFilePath(filePath string) string {
	originalPath := filePath
	i := 0
	for {
		if !FileExist(originalPath) {
			break
		}else {
			i += 1
			var builder strings.Builder
			builder.WriteString(GetFilePath(filePath) + "/")
			builder.WriteString(GetFileName(filePath))
			builder.WriteString("(" + strconv.Itoa(i) + ")")
			builder.WriteString(GetFileSuffixString(filePath))
			originalPath = builder.String()
		}
	}

	return originalPath
}