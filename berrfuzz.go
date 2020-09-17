package main

//https://www.w3schools.com/charsets/ref_html_utf8.asp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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

// SetupLogger sets up local log file for fuzzing output
func SetupLogger() {
	logFile, err := os.OpenFile("BerrFuzz-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//defer logFile.Close()
	logFile.WriteString("\n")
	log.SetOutput(logFile)
	log.Printf("-=-==-===-====Start of New Fuzz Test")
}

//Generator()
// --- This will generate a random file, which will then be output to be used by the fuzzer
//Mutator() This will mutate the existing file
// --- Should eventually develop into something that allows us to select what to mutate

func main() {
	fmt.Println("-=BerrFuzz")

	//Add options to delete local log file
	SetupLogger()

	fileBytes, err := ReadFileBytes("test.txt")
	if err != nil {
		fmt.Println("Error reading files")
	}
	fmt.Println(string(fileBytes))

	// Payload can be an optional input
	payload := "..\\"
	totalNum := 2000

	for i := 0; i < totalNum; i++ {
		payload += "..\\✁"
	}
	fmt.Printf(payload)

	//! This will be a command line input
	targetProgram := "Notepad"

	cmd := exec.Command(targetProgram, payload)
	log.Printf("Running command")
	//err = cmd.Run()

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("-=-=####### Possible Crash:")
		log.Printf("-=-=####### \tErr: %s", err.Error())
		log.Printf("-=-=#######\t\tPayload @ Possible Crash: \"%s\"", payload)
	}

	fmt.Printf("Output of program: %s", string(output))
	log.Print("Done running command")

	//TODO: Random byte + character generations
	//TODO: Being able to choose certain character sets
	//TODO: Integrate known bad strings
	//TODO: Possibly integrate search for known bad functions
}
