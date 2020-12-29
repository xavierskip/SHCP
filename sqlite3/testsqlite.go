package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/05.3.html
func main() {
	var ip string
	var mac string
	var name string
	var comment string
	var time string
	var lastdate string

	tablename := os.Args[1]
	queryip := os.Args[2]
	fmt.Println("os.Args", tablename, queryip)

	db, err := sql.Open("sqlite3", "./network.db")
	checkErr(err)

	//查询数据
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE ip=?;", tablename)
	stmt, err := db.Prepare(sql)
	checkErr(err)
	defer stmt.Close()
	err = stmt.QueryRow(queryip).Scan(&ip, &mac, &name, &comment, &time, &lastdate)
	checkErr(err)
	fmt.Println(ip, mac, name, comment, time, lastdate)

	db.Close()
	hw, _ := net.ParseMAC(mac)
	address := fmt.Sprintf("%s", hw)
	fmt.Printf("MAC: %s\n", address)

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
