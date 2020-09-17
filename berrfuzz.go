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

//Generator()
// --- This will generate a random file, which will then be output to be used by the fuzzer
//Mutator() This will mutate the existing file
// --- Should eventually develop into something that allows us to select what to mutate

func main() {
	fmt.Println("-=BerrFuzz")

	fileBytes, err := ReadFileBytes("test.txt")
	if err != nil {
		fmt.Println("Error reading files")
	}
	fmt.Println(string(fileBytes))

	temp := "A"

	for i := 0; i < 20; i++ {
		temp += "A"
		fmt.Printf("%s\n", temp)
	}
}
