package writefile

import (
	"fmt"
	"os"
	"time"
)

func IsValidFile(file string) bool {

	info, err := os.Stat(file)

	if os.IsNotExist(err) {
		return false
	} else if info.IsDir() {
		return false
	} else {
		return true
	}
}

func writeData(fname string, data []byte) {
	if !IsValidFile(fname) {
		ef, err := os.Create(fname)
		if err != nil {
			fmt.Printf("Unable to create file: %v", fname)
			return
		}
		ef.Close()
	}

	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening file: %v/n", fname)
		return
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		fmt.Printf("Error writing to file: %v/n", fname)
		return
	}
}

func WriteLog(fname string, text string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	l := "\r" + string(currentTime) + " " + text + "\n"
	fmt.Printf(l)
	data := []byte(l)
	writeData(fname, data)
}

func WriteFile(fname string, data []byte) {
	//fmt.Println("Writing file to: " + fname)
	writeData(fname, data)
}
