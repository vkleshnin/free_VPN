package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Host struct {
	Score          string
	Speed          string
	Uptime         string
	TotalUsers     string
	TotalTraffic   string
	HostName       string
	Ip             string
	CountryLong    string
	CountryShort   string
	LogType        string
	Operator       string
	Message        string
	Ping           string
	NumVpnSessions string
}

func (ObjectHost *Host) getInformation() {
	fmt.Printf("%-20s %20s\n", "Host Name:", ObjectHost.HostName)
	fmt.Printf("%-20s %20s\n", "IP:", ObjectHost.Ip)
	fmt.Printf("%-20s %20s\n", "Score:", ObjectHost.Score)
	fmt.Printf("%-20s %20s\n", "Ping:", ObjectHost.Ping)
	fmt.Printf("%-20s %20s\n", "Speed:", ObjectHost.Speed)
	fmt.Printf("%-20s %20s\n", "Country Long:", ObjectHost.CountryLong)
	fmt.Printf("%-20s %20s\n", "Country Short:", ObjectHost.CountryShort)
	fmt.Printf("%-20s %20s\n", "Num Vpn Sessions:", ObjectHost.NumVpnSessions)
	fmt.Printf("%-20s %20s\n", "Uptime:", ObjectHost.Uptime)
	fmt.Printf("%-20s %20s\n", "Total Users:", ObjectHost.TotalUsers)
	fmt.Printf("%-20s %20s\n", "Total Traffic:", ObjectHost.TotalTraffic)
	fmt.Printf("%-20s %20s\n", "Log Type:", ObjectHost.LogType)
	fmt.Printf("%-20s %20s\n", "Operator:", ObjectHost.Operator)
	fmt.Printf("%-20s %20s\n", "Message:", ObjectHost.Message)
	fmt.Printf("-----------------------------------------\n")
}

func (ObjectHost *Host) isFineVPN() int {
	if !strings.HasPrefix(ObjectHost.HostName, "public") {
		return 1
	}
	temp1, err := strconv.Atoi(ObjectHost.Speed)
	if err != nil {
		return 1
	}
	if float64(temp1)/1024/1024 > 100 {
		ObjectHost.Speed = fmt.Sprintf("%.1f Mbps", float64(temp1)/1024/1024)
	} else {
		return 1
	}
	temp2, err := strconv.Atoi(ObjectHost.Ping)
	if temp2 > 15 {
		return 1
	}
	return 0
}

func downloadFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal()
		}
	}(resp.Body)
	out, err := os.Create("ip.csv")
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal()
		}
	}(out)
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	if _, err := os.Stat("ip.csv"); os.IsNotExist(err) {
		err := downloadFile("https://www.vpngate.net/api/iphone")
		if err != nil {
			panic(interface{}("file is not download"))
		}
	}
	csvFile, err := os.OpenFile("ip.csv", os.O_RDWR, 0600)
	if err != nil {
		panic(interface{}("file is not open"))
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			panic(interface{}("file is not closed"))
		}
	}(csvFile)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comment = '*'
	var ObjectHost []Host
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		ObjectHost = append(ObjectHost, Host{
			HostName:       line[0],
			Ip:             line[1],
			Score:          line[2],
			Ping:           line[3],
			Speed:          line[4],
			CountryLong:    line[5],
			CountryShort:   line[6],
			NumVpnSessions: line[7],
			Uptime:         line[8],
			TotalUsers:     line[9],
			TotalTraffic:   line[10],
			LogType:        line[11],
			Operator:       line[12],
			Message:        line[13],
		})
	}
	for i := range ObjectHost {
		if ObjectHost[i].isFineVPN() == 0 {
			ObjectHost[i].getInformation()
		}
	}
}
