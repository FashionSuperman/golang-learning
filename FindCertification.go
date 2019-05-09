package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const ERRTEXT = "对不起，没有查询到相关证书数据或查询条件有误，请重新确认后再查询"

//const ServiceUrl = "http://124.128.231.226:8888/sdostair/sdosta/searchResult/getresultWeb?type=web&aac002=%s&aac003=%s"
const ServiceUrl = "http://124.128.231.226:8888/sdostair/sdosta/searchResult/yanZsWeb?type=web&aac002=&aac003=%s&bzb178=%s"
const SheetName = "Sheet1"
const AxisColumn = "G"
const NameColumn = 1
const IdColumn = 2
const CerIdColumn = 7

func main() {
	//解析文件参数
	var filename string
	filename = os.Args[1]
	fmt.Println("fileName is ", filename)
	//filename = "./2.xlsx"
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "where is the file? %v\n", err)
		return
	}
	dealFile(f)
}

func dealFile(file *excelize.File) {
	idToNameMap := make(map[string]string)
	idToCertifications := make(map[string]string)
	rows, err := file.GetRows(SheetName)
	if err != nil {
		fmt.Println("sheet is not valid")
		return
	}
	for _, row := range rows {
		//for _, colCell := range row {
		//	fmt.Print(colCell, "\t")
		//}
		//fmt.Println()

		if len(row) < 3 {
			continue
		}
		cellName := row[NameColumn]
		cellId := row[IdColumn]
		cellCerId := row[CerIdColumn]
		if len(cellId) != 18 {
			fmt.Println("身份证编码错误", cellName, "===", cellId)
			continue
		}
		//fmt.Print(cellName, "\t" , cellId, "\t")
		//fmt.Println()
		idToNameMap[cellId] = cellName + "," + cellCerId
	}
	for key, value := range idToNameMap {
		fmt.Println("key is ", key, " value is ", value)
	}
	fmt.Println(len(idToNameMap))

	//
	for key, value := range idToNameMap {
		tempArr := strings.Split(value, ",")
		name := tempArr[0]
		cerId := tempArr[1]
		urlQuery := strings.Replace(ServiceUrl, "%s", url.QueryEscape(name), 1)
		urlQuery = strings.Replace(urlQuery, "%s", cerId, 1)
		respose, err := http.Get(urlQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "请求服务失败 %v\n", err)
			continue
		}
		if 200 != respose.StatusCode {
			fmt.Println("状态码错误:", respose.StatusCode, " urlQuery is:", urlQuery)
			continue
		}

		//data,e := ioutil.ReadAll(respose.Body)
		//if e != nil {
		//	fmt.Fprintf(os.Stderr, "处理数据失败 %v\n", err)
		//	continue
		//}
		//tempBodyStr := string(data)
		//if strings.Contains(tempBodyStr,"对不起，没有查询到相关证书数据或查询条件有误，请重新确认后再查询") {
		//	fmt.Println("没有查询到数据-", value)
		//	continue
		//}

		doc, e := goquery.NewDocumentFromReader(respose.Body)
		respose.Body.Close()
		if e != nil {
			fmt.Fprintf(os.Stderr, "转换doc失败 %v\n", e)
			continue
		}

		sel := doc.Find(".un_sorry")
		if sel.Length() != 0 && strings.Contains(sel.Find("span").Text(), ERRTEXT) {
			fmt.Println("没有查询到数据-", value)
			continue
		}

		doc.Find(".lwcx-table").Each(func(i int, selection *goquery.Selection) {
			certificateName := selection.Find("tbody tr:nth-child(2) td:nth-child(2)").Text()
			title := selection.Find("tbody tr:nth-child(2) td:nth-child(1)").Text()
			if title != "鉴定工种" {
				return
			}
			idToCertifications[key] = certificateName
		})

	}

	for index, row := range rows {
		if len(row) < 3 {
			continue
		}
		cellId := row[IdColumn]
		if len(cellId) != 18 {
			continue
		}
		axis := AxisColumn + strconv.Itoa(index+1)
		for id, certification := range idToCertifications {
			if id == cellId {
				file.SetCellValue(SheetName, axis, certification)
			}
		}
	}

	file.Save()
}
