package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetDuration(createionTime, completionTime time.Time) string {
	duration := completionTime.Sub(createionTime)

	durationTime := ""
	//fmt.Println(int(duration.Seconds()), int(duration.Minutes()), int(duration.Seconds()) % 60)
	if duration.Hours() >= 240 {
		durationTime = strconv.Itoa(int(duration.Hours()/24)) + "d"
	}else if duration.Hours() >= 24 {
		durationTime = strconv.Itoa(int(duration.Hours()/24)) + "d" + strconv.Itoa(int(duration.Hours())%24) + "h"
	} else if duration.Hours() >= 10 {
		durationTime = strconv.Itoa(int(duration.Hours())) + "h"
	} else if duration.Hours() >= 1 {
		durationTime = strconv.Itoa(int(duration.Hours())) + "h" + strconv.Itoa(int(duration.Minutes())%60) + "m"
	} else if duration.Minutes() >= 10 {
		durationTime = strconv.Itoa(int(duration.Minutes())) + "m"
	} else if duration.Minutes() >= 1 {
		durationTime = strconv.Itoa(int(duration.Minutes())) + "m" + strconv.Itoa(int(duration.Seconds())%60) + "s"
	} else {
		durationTime = strconv.Itoa(int(duration.Seconds())) + "s"
	}
	return durationTime
}

func GetAge(createionTime time.Time) string {
	duration := time.Since(createionTime)
	age := ""
	//fmt.Println(int(duration.Seconds()), int(duration.Minutes()), int(duration.Seconds()) % 60)
	if duration.Hours() >= 240 {
		age = strconv.Itoa(int(duration.Hours()/24)) + "d"
	}else if duration.Hours() >= 24 {
		age = strconv.Itoa(int(duration.Hours()/24)) + "d" + strconv.Itoa(int(duration.Hours())%24) + "h"
	} else if duration.Hours() >= 10 {
		age = strconv.Itoa(int(duration.Hours())) + "h"
	} else if duration.Hours() >= 1 {
		age = strconv.Itoa(int(duration.Hours())) + "h" + strconv.Itoa(int(duration.Minutes())%60) + "m"
	} else if duration.Minutes() >= 10 {
		age = strconv.Itoa(int(duration.Minutes())) + "m"
	} else if duration.Minutes() >= 1 {
		age = strconv.Itoa(int(duration.Minutes())) + "m" + strconv.Itoa(int(duration.Seconds())%60) + "s"
	} else {
		age = strconv.Itoa(int(duration.Seconds())) + "s"
	}
	return age
}
func GetFileNameList() []string{
	fileOrDirname, _ := filepath.Abs(Option_file)
	filenameList := []string{}

	fi, err := os.Stat(fileOrDirname)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
	// do directory stuff

		files, err := ioutil.ReadDir(fileOrDirname)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if err != nil {
				fmt.Println(err)
			}
			if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml"{
				filenameList = append(filenameList, f.Name())
			}
		}
	case mode.IsRegular():
		// do file stuff

		filenameList = append(filenameList, fileOrDirname)
	}
	return filenameList
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}