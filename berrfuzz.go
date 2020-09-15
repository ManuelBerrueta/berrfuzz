package main

import (
	"bufio"
	"fmt"
	"os"
)

// ReadFileBytes reads a file and returns bytes
func ReadFileBytes(fileName string) ([]byte, error) {
	inFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("inFile Read Error")
		return nil, err
	}
	defer inFile.Close()

	inFileStats, errStats := inFile.Stat()
	if errStats != nil {
		fmt.Println("File Stats Read Error")
		return nil, errStats
	}

	fileSize := inFileStats.Size()
	fileBytes := make([]byte, fileSize)

	inFileBuffer := bufio.NewReader(inFile)
	_, err = inFileBuffer.Read(fileBytes)

	return fileBytes, err
}

func main() {
	fmt.Println("BerrFuzz")

	fileBytes, err := ReadFileBytes("test.txt")
	if err != nil {
		fmt.Println("Error reading files")
	}
	fmt.Println(string(fileBytes))
	//data, err := ioutil.ReadFile("test.txt")
	//if err != nil {
	//	fmt.Println("File reading error", err)
	//	return
	//}
	//fmt.Println("Contents of file:", string(data))

}
