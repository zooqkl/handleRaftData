package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"io"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

const (
	BLOCKIO  = "blockIO"
	TLS      = "tls"
	MEMORY   = "memory"
	NETSPEED = "netSpeed"
	CPU      = "cpu"

	SHEETNAME_ONCE_TRAN_TIME = "onceTranTime"
	SHEETNAME_TPS            = "tps"

	XLSX_PATH = "./data/xlsx/"
)

type SourcePath struct {
	testType       string
	sourceFilePath []string
}

func handleRaftTestData(sourcePathMap map[string][]SourcePath) {

	for testType, sourcePaths := range sourcePathMap {
		xlsdata := make(map[string][][]string)
		xlsdata_sheet_once := make([][]string, 0)
		xlsdata_sheet_tps := make([][]string, 0)
		for _, sourcePath := range sourcePaths {
			for _, filePath := range sourcePath.sourceFilePath {
				xlsdata_column, err := readRaftXlsxFile(filePath)
				if err != nil {
					panic(err)
				}
				if strings.Contains(filePath, SHEETNAME_ONCE_TRAN_TIME) {
					xlsdata_sheet_once = append(xlsdata_sheet_once, xlsdata_column)
				} else {
					xlsdata_sheet_tps = append(xlsdata_sheet_tps, xlsdata_column)
				}
			}
			xlsdata[SHEETNAME_ONCE_TRAN_TIME] = xlsdata_sheet_once
			xlsdata[SHEETNAME_TPS] = xlsdata_sheet_tps
		}
		createRaftXlsxFile(testType+" impact on performance", xlsdata)
	}
}

func readRaftXlsxFile(xlsPath string) ([]string, error) {
	resultData := make([]string, 0)
	str := strings.Split(strings.Split(xlsPath, "/resultExcel")[0], "/")

	fi, err := os.Open(xlsPath)
	if err != nil {
		return resultData, fmt.Errorf("Error: %s\n", err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		dataBytes, _, c := br.ReadLine()
		results := strings.Split(string(dataBytes), ",")
		resultData = append(resultData, results[len(results)-1])
		if c == io.EOF {
			fmt.Printf("read file %s finish!\n", str[3:])
			break
		}
	}
	resultData[0] = str[len(str)-1]
	return resultData, nil
}

func createRaftXlsxFile(fileName string, data map[string][][]string) {

	xlsx := excelize.NewFile()

	for sheetName, sheetData := range data {
		index := xlsx.NewSheet(sheetName)
		for columnIndex, columnData := range sheetData {
			for rowIndex, data := range columnData {
				letter, err := intConverLetter(columnIndex + 1)
				if err != nil {
					panic(err)
				}
				axis := letter + strconv.Itoa(rowIndex+1)
				xlsx.SetCellStr(sheetName, axis, data)
			}
		}
		xlsx.SetActiveSheet(index)
	}
	xlsx.DeleteSheet("Sheet1")

	isExists, err := pathExists(XLSX_PATH)
	if err != nil || !isExists {
		err = os.Mkdir(XLSX_PATH, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	err = xlsx.SaveAs(XLSX_PATH + fileName + ".xlsx")
	if err != nil {
		panic(err)
	}
	fmt.Printf("xlsFile[%s] create Success!\n", fileName)
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func intConverLetter(num int) (string, error) {
	switch num {
	case 1:
		return "A", nil
	case 2:
		return "B", nil
	case 3:
		return "C", nil
	case 4:
		return "D", nil
	case 5:
		return "E", nil
	case 6:
		return "F", nil
	case 7:
		return "G", nil
	case 8:
		return "H", nil
	case 9:
		return "I", nil
	case 10:
		return "J", nil
	case 11:
		return "K", nil
	case 12:
		return "L", nil
	case 13:
		return "M", nil
	case 14:
		return "N", nil
	case 15:
		return "O", nil
	case 16:
		return "P", nil
	case 17:
		return "Q", nil
	case 18:
		return "R", nil
	case 19:
		return "S", nil
	case 20:
		return "T", nil
	case 21:
		return "U", nil
	case 22:
		return "V", nil
	case 23:
		return "W", nil
	case 24:
		return "X", nil
	case 25:
		return "Y", nil
	case 26:
		return "Z", nil
	default:
		return "", fmt.Errorf("Incorrect parameter Error!  parameter format is [1,26] !")
	}
}

func sourceFilePathFormat() map[string][]SourcePath {
	sourcePathMap := make(map[string][]SourcePath)
	testTypes := make(map[string][]string)

	testTypes[BLOCKIO] = []string{"500mb", "10mb", "5mb", "1mb", "512kb"}
	testTypes[CPU] = []string{"4core", "3core", "2core", "1core", "0.5core", "0.1core"}
	testTypes[MEMORY] = []string{"8g", "4g", "2g", "1g", "512m", "256m"}
	testTypes[NETSPEED] = []string{"100mb:s", "10mb:s", "5mb:s", "1mb:s", "512kb:s", "100kb:s",}
	testTypes[TLS] = []string{"notTls", "tls"}

	path1 := "./data/性能测试/oldPath/resultExcel/onceTranTime.xls"
	path2 := "./data/性能测试/oldPath/resultExcel/peer0.org1.example.com.xls"

	for testType, typeDetails := range testTypes {
		sourcePaths := make([]SourcePath, 0)
		for _, typevalue := range typeDetails {
			paths := make([]string, 0)
			newPath := testType + "/" + typevalue
			path_1 := strings.Replace(path1, "oldPath", newPath, 1)
			path_2 := strings.Replace(path2, "oldPath", newPath, 1)
			paths = append(append(paths, path_1), path_2)
			sourcePaths = append(sourcePaths, SourcePath{sourceFilePath: paths, testType: typevalue})
		}
		sourcePathMap[testType] = sourcePaths
	}
	return sourcePathMap
}

func main() {
	data := sourceFilePathFormat()
	handleRaftTestData(data)
}
