package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Event struct {
	Name     string `json:"name"`
	Link     string `json:"link"`
	Type     string `json:"type"`
	Bmks     string `json:"bmks"`
	Bmjz     string `json:"bmjz"`
	Bsks     string `json:"bsks"`
	Bsjs     string `json:"bsjs"`
	Readmore string `json:"readmore"`
	ID       int    `json:"id"`
	Status   int    `json:"status"`
}

type MatchList struct {
	competitionID                         int
	competitionName, registrationDeadline string
}

type Result struct {
	Events []Event `json:"result"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
}

type ApiResponse struct {
	Success bool   `json:"success"`
	Data    Result `json:"data"`
	Msg     string `json:"msg"`
}

func Matchs() string {
	if UpdataTable() {
		Getdata()
		fmt.Println("比赛信息更新")
	}

	return GetMatchs()
}

func Getdata() string {

	var ApiData ApiResponse
	data := map[string]interface{}{}
	payload, err := json.Marshal(data)

	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return ""
	}
	// 发送 POST 请求
	response, err := http.Post("https://www.su-sanha.cn/api/events/list", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("发送请求时发生错误:", err)
		return ""
	}
	defer response.Body.Close()

	// 读取响应的 JSON 数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应时发生错误:", err)
		return ""
	}

	// 解析 JSON 数据

	err = json.Unmarshal(body, &ApiData)
	if err != nil {
		fmt.Println("解析 JSON 时发生错误:", err)
		return ""
	}

	mysql_data(ApiData)

	return "查询成功"

}

func mysql_data(ApiData ApiResponse) string {
	layout := "2006年01月02日 15:04"
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, data := range ApiData.Data.Events {
		RegistrationStart, _ := time.Parse(layout, data.Bmks)
		RegistrationDeadline, _ := time.Parse(layout, data.Bmjz)
		CompetitionStart, _ := time.Parse(layout, data.Bsks)
		CompetitionEnd, _ := time.Parse(layout, data.Bsjs)

		// 准备 INSERT 语句
		query := `
		INSERT INTO Competition (
			CompetitionName,
			CompetitionLink,
			CompetitionType,
			RegistrationStart,
			RegistrationDeadline,
			CompetitionStart,
			CompetitionEnd,
			AdditionalInfo
		)
		SELECT
			?, ?, ?, ?, ?, ?, ?, ?
		FROM
			DUAL
		WHERE
			NOT EXISTS (
				SELECT 1
				FROM Competition
				WHERE CompetitionName = ?
			)
	`

		// 执行 INSERT 语句
		_, err = db.Exec(query,
			data.Name,
			data.Link,
			data.Type,
			RegistrationStart,
			RegistrationDeadline,
			CompetitionStart,
			CompetitionEnd,
			data.Readmore,
			data.Name,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	return "成功!"

}

func GetMatchs() string {
	layout := "2006-01-02 15:04:05"
	var matchs string
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query2 := `
		SELECT CompetitionID, CompetitionName, RegistrationDeadline
		FROM Competition
		ORDER BY RegistrationDeadline
	`

	// 执行查询语句
	rows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var competition MatchList

		err := rows.Scan(&competition.competitionID, &competition.competitionName, &competition.registrationDeadline)
		if err != nil {
			log.Fatal(err)
		}
		expiration, err := time.Parse(layout, competition.registrationDeadline)

		if err != nil {
			fmt.Println("Error parsing time:", err)
		}

		currentTime := time.Now()

		if !expiration.Before(currentTime) {

			matchs = matchs + fmt.Sprintf("ID:%v\n比赛名称%v\n报名截止时间%v\n----------------\n", competition.competitionID, competition.competitionName, competition.registrationDeadline)
		}
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return matchs
}

func UpdataTable() bool {

	db, err := sql.Open("mysql", "kali:kali@tcp(localhost:3306)/boot?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 执行查询
	rows, err := db.Query("SELECT TableLastUpdated FROM Competition ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	flag := 0
	// 处理查询结果
	for rows.Next() {
		flag = 1
		var tableLastUpdated sql.NullTime
		err := rows.Scan(&tableLastUpdated)
		if err != nil {
			log.Fatal(err)
		}
		// 比较时间
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		if tableLastUpdated.Time.Before(oneHourAgo) {
			return true
		}
	}

	// 检查错误
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	if flag == 0 {
		return true
	}
	return false
}

func Match(num int) string {

	db, err := sql.Open("mysql", "kali:kali@tcp(localhost:3306)/boot?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 需要查询的 CompetitionID
	var competitionID = num

	// 执行查询
	var comp Event
	query := `SELECT  CompetitionName, CompetitionLink, CompetitionType, 
		RegistrationStart, RegistrationDeadline, CompetitionStart, CompetitionEnd, AdditionalInfo 
		FROM Competition WHERE CompetitionID = ?`

	err = db.QueryRow(query, competitionID).Scan(&comp.Name, &comp.Link,
		&comp.Type, &comp.Bmks, &comp.Bmjz, &comp.Bsks,
		&comp.Bsjs, &comp.Readmore)

	if err != nil {
		log.Fatal(err)
	}

	MathcData := fmt.Sprintf("比赛名称：%v\n比赛类型：%v\n报名开始：%s\n报名截止：%s\n比赛开始：%s\n比赛结束：%s\n其他说明：%v", comp.Name, comp.Type, TimePa(comp.Bmks), TimePa(comp.Bmjz), TimePa(comp.Bsks), TimePa(comp.Bsjs), comp.Readmore)
	return MathcData
}

func TimePa(timeData string) string {
	// 解析日期字符串
	originalTime, err := time.Parse(time.RFC3339, timeData)
	if err != nil {
		fmt.Println("日期解析失败:", err)
		return ""
	}

	// 格式化为目标日期字符串
	targetFormat := "2006年01月02日 15:04:05"
	targetDateString := originalTime.Format(targetFormat)

	return targetDateString
}
