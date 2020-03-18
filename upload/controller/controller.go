package uploadController

import (
	"net/http"
	"log"
	"os"
	"io"
	"mime/multipart"
	"errors"
	"encoding/csv"
	"github.com/axgle/mahonia"
	"github.com/tealeg/xlsx"
	"lab-golang/tool/yaml"
	"lab-golang/tool/tool"
	"path/filepath"
	"strings"
)

type FileCell struct {
	column int
	value string
}

type FileRow struct {
	row int
	cells []FileCell
}

type CSVFile struct {
	fileName string
	path string
	data []FileRow
}

type LabFileType string

const (
	LabFileTypeUnknown LabFileType = ""
	LabFileTypeTxt LabFileType = "txt"
	LabFileTypeXlsx LabFileType = "xlsx"
	LabFileTypeCSV LabFileType = "csv"

	FastDelivery = "快速发货"
	NoneAfterSalesStatus = "无售后或售后取消"
	ToBeDelivered = "待发货"
	NaturalGas = "天然气"
	LiquefiedGas = "液化气"
	WordValve = "阀"
	WordGeniusOnsiteInstallation = "全国联保上门安装"
	WordGenius = "全国联保"
	WordOnsiteInstallation = "上门安装"
)

func SaveFile(w http.ResponseWriter, r *http.Request, fileId string) (string, string, error) {

	r.ParseForm()
	//把上传的文件存储在内存和临时文件中
	err := r.ParseMultipartForm(32<<20)

	if err != nil {
		log.Printf("解析文件错误: %v", err.Error())
		return "", "", err
	}

	xlsxFile, xlsxHandler, xlsxErr := r.FormFile(fileId)
	if xlsxErr != nil {
		log.Printf("xlsx文件错误: %v", xlsxErr.Error())
		return "", "", xlsxErr
	}
	defer xlsxFile.Close()

	targetFile := r.PostFormValue("fileName")

	filePath, fileErr := saveFileToLocalPath(xlsxHandler.Filename, xlsxFile, xlsxHandler)

	if fileErr != nil {
		log.Printf("错误原因: %v", fileErr.Error())
		return "", xlsxHandler.Filename, fileErr
	}

	log.Printf("文件名称:%v \n", xlsxHandler.Filename)

	return filePath, targetFile, nil
}

func saveFileToLocalPath(fileName string, file multipart.File, header *multipart.FileHeader) (string, error) {
    //创建上传的目的文件
    filePath := "./upload/file/" + fileName

    absoultePath, err := filepath.Abs(filePath)
    if err != nil {
    	return "", err
    }

    absoultePath = tool.NewFilePath(absoultePath)
    log.Printf("csv文件地址: %v \n", absoultePath)

    f, err := os.OpenFile(absoultePath, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0666)

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

    return absoultePath, nil
}

func ReadCSVFile(filePath string, fileName string) (*CSVFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("打开文件错误: %v\n", err.Error())
		return nil, err
	}

	// decoder := mahonia.NewDecoder("gbk")
	// reader := csv.NewReader(decoder.NewReader(file))
	reader := csv.NewReader(file)
	if reader == nil {
		log.Printf("初始化读取失败\n")
		return nil, errors.New("文件初始化失败")
	}

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("读取数据失败: %v \n", err.Error())
		return  nil, err
	}

	colNum := len(records[0])
	recordNum := len(records)

	if recordNum <= 0 {
		return nil, errors.New("文件为空")
	}

	data := make([]FileRow, 0)
	for i := 0; i < recordNum; i++ {
		row := FileRow{row: i}
		rowValues := make([]FileCell, 0)
		for k := 0; k < colNum; k++ {
			cell := FileCell{column:k}
			cell.value = records[i][k]
			rowValues = append(rowValues, cell)
		}
		row.cells = rowValues
		data = append(data, row)
	}

	csvFile := &CSVFile{fileName: fileName, path : filePath, data: data}

	return csvFile, nil
}

func SaveDataToXlsxFile(data *CSVFile) (string, error) {
	if len(data.data) == 0 {
		return "", errors.New("数据为空，创建文件失败")
	}


	newFile := xlsx.NewFile()
	newSheet, err := newFile.AddSheet("newSheet")
	if err != nil {
		log.Println("新建excel文件清单失败")
		return "", errors.New("新建清单失败:" + err.Error())
	}

	for _, row := range data.data {
		newRow := newSheet.AddRow()
		for _, cell := range row.cells {
			newCell := newRow.AddCell()
			newCell.SetValue(cell.value)
			newCell.SetStyle(xlsx.NewStyle())
		}
	}

	if len(newSheet.Rows) == 0 {
		return "", errors.New("匹配数据失败")
	}

	reader, err := yamlReader.Instance()
	if err != nil {
		return "", err
	}


	newFilePath := reader.Configure.Xlsx.SavedDirctory + tool.GetFileName(data.fileName) + tool.AppendFileSuffix(data.fileName, "xlsx")

	newAbsoulteFilePath, err := filepath.Abs(newFilePath)
	if err != nil {
		return "", err
	}

	newAbsoulteFilePath = tool.NewFilePath(newAbsoulteFilePath)

	newErr := newFile.Save(newAbsoulteFilePath)

	if err != nil {
		log.Printf("创建文件失败[%v]\n", newErr.Error())
		return "", errors.New("创建文件" + newAbsoulteFilePath + "失败:" + err.Error())
	}


	return newAbsoulteFilePath, nil
}

/*
	过滤筛选数据，并且生成xlsx文件
*/
func CreateNewShippingOrder(data *CSVFile) (string, error) {
	if len(data.data) == 0 {
		return "", errors.New("数据为空，创建文件失败")
	}


	newFile := xlsx.NewFile()
	newSheet, err := newFile.AddSheet("newSheet")
	if err != nil {
		log.Println("新建excel文件清单失败")
		return "", errors.New("新建清单失败:" + err.Error())
	}

	/* sheet 标题行 */
	titles := []string{"订单号", "收货人", "手机", "地址", "数量", "物品名", "快递单号列"}
	newRow := newSheet.AddRow()
	log.Printf("当前标题: %v\n", strings.Join(titles, " "))
	for _ , title := range titles {
		headerCell := newRow.AddCell()
		headerCell.SetValue(title)
		headerCell.SetStyle(xlsx.NewStyle())
	}

	for _, row := range data.data {
		log.Printf("当前内容: %v", row.cells)
		/*
		AN 商家备注 过滤掉"快速发货"
		AO 售后状态 保留 "无售后或售后取消"
		C  订单状态 保留 "待发货"
		*/
		if row.cells[2].value == ToBeDelivered && row.cells[40].value == NoneAfterSalesStatus && !strings.Contains(row.cells[39].value, FastDelivery) {
			newRow := newSheet.AddRow()
			// 订单号 + 收货人 + 手机 + 地址 + 数量 + 物品名 + 快递单号列

			orderCell := newRow.AddCell()
			orderCell.SetValue(row.cells[1].value)
			orderCell.SetStyle(xlsx.NewStyle())

			userCell := newRow.AddCell()
			userCell.SetValue(row.cells[14].value)
			userCell.SetStyle(xlsx.NewStyle())

			phoneCell := newRow.AddCell()
			phoneCell.SetValue(row.cells[15].value)
			phoneCell.SetStyle(xlsx.NewStyle())

			addressCell := newRow.AddCell()
			var builder strings.Builder
			builder.WriteString(row.cells[17].value)
			builder.WriteString(row.cells[18].value)
			builder.WriteString(row.cells[19].value)
			builder.WriteString(row.cells[20].value)
			addressCell.SetValue(builder.String())
			addressCell.SetStyle(xlsx.NewStyle())

			numberCell := newRow.AddCell()
			numberCell.SetValue(row.cells[11].value)
			numberCell.SetStyle(xlsx.NewStyle())


			goods_sku := row.cells[30].value
			suffixValue := ""
			copyCurerntRow := false
			// 商品规格
			kindsValue := row.cells[27].value
			if strings.Contains(kindsValue, NaturalGas) || strings.Contains(kindsValue, LiquefiedGas) {
				if strings.Contains(kindsValue, WordValve) {
					// 复制一份当前行
					suffixValue = WordValve
					copyCurerntRow = true
				}
				if strings.Contains(kindsValue, NaturalGas) {
					suffixValue = NaturalGas
				}else if strings.Contains(kindsValue, LiquefiedGas) {
					suffixValue = LiquefiedGas
				}
			}else {
				// 含 "全国联保"
				if strings.Contains(kindsValue, WordGenius) {
					if strings.Contains(kindsValue, WordGeniusOnsiteInstallation) {
						suffixValue = WordGeniusOnsiteInstallation
					}else {
						suffixValue = WordGenius
					}
				}
			}

			goodsValue := goods_sku + " " + suffixValue
			if len(goods_sku) == 0 {
				suffixValue = strings.ReplaceAll(suffixValue, " ", "")
				if len(suffixValue) == 0 {
					goodsValue = ""
				}else {
					goodsValue = suffixValue;
				}
			}else {
				suffixValue = strings.ReplaceAll(suffixValue, " ", "")
				if len(suffixValue) == 0 {
					goodsValue = goods_sku
				}
			}



			goodsNameCell := newRow.AddCell()
			goodsNameCell.SetValue(goodsValue)
			goodsNameCell.SetStyle(xlsx.NewStyle())

			expressCell := newRow.AddCell()
			expressCell.SetStyle(xlsx.NewStyle())
			log.Printf("当前匹配到的数据: %v | %v | %v | %v | %v | %v | %v\n", orderCell.String(), userCell.String(), phoneCell.String(), addressCell.String(), numberCell.String(), goodsNameCell.String(), expressCell.String())

			if copyCurerntRow {
				copyCurerntRow = false
				log.Printf("复制匹配的数据: ")
				newRow := newSheet.AddRow()
				// 订单号 + 收货人 + 手机 + 地址 + 数量 + 物品名 + 快递单号列

				orderCell := newRow.AddCell()
				orderCell.SetValue(row.cells[1].value)
				orderCell.SetStyle(xlsx.NewStyle())

				userCell := newRow.AddCell()
				userCell.SetValue(row.cells[14].value)
				userCell.SetStyle(xlsx.NewStyle())

				phoneCell := newRow.AddCell()
				phoneCell.SetValue(row.cells[15].value)
				phoneCell.SetStyle(xlsx.NewStyle())

				addressCell := newRow.AddCell()
				var builder strings.Builder
				builder.WriteString(row.cells[17].value)
				builder.WriteString(row.cells[18].value)
				builder.WriteString(row.cells[19].value)
				builder.WriteString(row.cells[20].value)
				addressCell.SetValue(builder.String())
				addressCell.SetStyle(xlsx.NewStyle())

				numberCell := newRow.AddCell()
				numberCell.SetValue(row.cells[11].value)
				numberCell.SetStyle(xlsx.NewStyle())

				goodsNameCell := newRow.AddCell()
				goodsNameCell.SetValue(WordValve + row.cells[11].value + "个")
				goodsNameCell.SetStyle(xlsx.NewStyle())

				expressCell := newRow.AddCell()
				expressCell.SetStyle(xlsx.NewStyle())

				log.Printf("当前匹配到的数据: %v | %v | %v | %v | %v | %v ｜ %v\n", orderCell.String(), userCell.String(), phoneCell.String(), addressCell.String(), numberCell.String(), goodsNameCell.String(), expressCell.String())
			}
		}
	}

	if len(newSheet.Rows) == 0 {
		return "", errors.New("匹配数据失败")
	}

	reader, err := yamlReader.Instance()
	if err != nil {
		return "", err
	}

	newFileName := tool.AppendFileSuffix(data.fileName, "xlsx")

	newFilePath := reader.Configure.Xlsx.SavedDirctory + newFileName

	newAbsoulteFilePath := tool.NewFilePath(newFilePath)

	newErr := newFile.Save(newAbsoulteFilePath)

	if err != nil {
		log.Printf("创建文件失败[%v]\n", newErr.Error())
		return "", errors.New("创建文件" + newAbsoulteFilePath + "失败:" + err.Error())
	}

	newFileName = tool.AppendFileSuffix(tool.GetFileName(newAbsoulteFilePath), "xlsx")

	return reader.Configure.Xlsx.DownloadFile + "/" + newFileName, nil
}
