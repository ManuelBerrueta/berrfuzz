package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"os"
)

// ReadFileBytes reads a file and returns bytes
func ReadFileBytes(fileName string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(fileName)
	return fileBytes, err
}

// MutateCorpus makes random changes to the bytes on the input file
//func MutateCorpus(originalBytes []byte) ([]byte, error) {
//
//	//TODO: Byte manipulation
//	//fileBytes := make([]byte, len(originalBytes))
//	mutatedBytes :=
//
//	return fileBytes,
//}

// CalculatePayloadSize takes the payloadSizePtr and figures out the size if percentage
func CalculatePayloadSize(payloadSizeFl float64, corpusFileBytes []byte) int {
	payloadSize := 0
	if payloadSizeFl != 0.0 {
		/* If payloadSizePtr is < 1.0 we should have an input file and we will
		   calculate the size as a percentage of that file size */
		if payloadSizeFl < 1.0 && payloadSizeFl > 0.0 {
			//* Check if a filename is passed
			if len(corpusFileBytes) == 0 {
				fmt.Println("Error: Passed percentage but no input file!")
				os.Exit(-1)
			} else {
				payloadSize = int(float64(len(corpusFileBytes)) * payloadSizeFl)
			}
		} else if payloadSizeFl >= 1.0 {
			payloadSize = int(payloadSizeFl)
		} else {
			fmt.Println("Payload size must be a positive number!")
			os.Exit(-1)
		}
	} else if len(corpusFileBytes) != 0 {
		payloadSize = mrand.Intn(len(corpusFileBytes))
	} else {
		payloadSize = mrand.Intn(0x1fffffe8)
	}
	return payloadSize
}

// CorpusIntake takes the file name for the input file, reads file as bytes to be used for fuzzing
func CorpusIntake(corpusName string) []byte {
	corpusFileBytes := make([]byte, 0)
	if corpusName != "" {
		var corpusErr error
		corpusFileBytes, corpusErr = ReadFileBytes(corpusName)
		if corpusErr != nil {
			fmt.Println("Error reading file: ", corpusName)
			os.Exit(-1)
		}
	}
	return corpusFileBytes
}

// OutputCorpus write the new payload output file
func OutputCorpus(fileName string, mutatedBytes []byte) bool {
	outFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("outFile Open/Create Error")
		return false
	}
	defer outFile.Close()

	numOfWrittenBytes, err := outFile.Write(mutatedBytes)
	if err != nil {
		fmt.Println("Corpus outFile Write Error")
	}

	if numOfWrittenBytes < len(mutatedBytes) {
		fmt.Println("numOfWrittenBytes < len(mutatedBytes)")
		return false
	}

	return true
}

// RandomByteGenerator generates random bytes of passed in size and returns []byte
func RandomByteGenerator(size int) []byte {
	builtBytes := make([]byte, size)

	_, err := rand.Read(builtBytes)

	if err != nil {
		fmt.Println("Error creating random bytes: ", err)
		os.Exit(-1)
	}

	return builtBytes
}

// RandomByteFileGenerator will create a file with random bytes
func RandomByteFileGenerator(size int, outFileName string) {
	builtBytes := make([]byte, size)
	_, err := rand.Read(builtBytes)

	if err != nil {
		fmt.Println("Error creating random bytes: ", err)
		os.Exit(-1)
	}

	outFile, err := os.OpenFile(outFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("outFile for RandomByteFileGenerator Failed to open")
		os.Exit(-1)
	}

	numByteWritten, err := outFile.Write(builtBytes)
	if err != nil {
		fmt.Println("outFile for RandomByteFileGenerator Failed to write bytes")
		os.Exit(-1)
	}

	if numByteWritten < len(builtBytes) {
		errorMessage := `Bytes written to outFile RandomByteFileGenerator
						less then length of random bytes to write`
		fmt.Println(errorMessage)
	}
}

// FileMutator is the start of a more complex function that will use different
//             techniques to flip bits/bytes
func FileMutator(fileBytes []byte, numOfByteToFlip int) []byte {
	//numOfByteToFlip := 10

	//C
	//randByte := mrand.Intn(0xFFFFFFFF)
	randByte := mrand.Intn(0xFF)

	for i := 0; i < numOfByteToFlip; i++ {
		fileBytes[mrand.Intn(len(fileBytes)-1)] = byte(randByte)
	}
	return fileBytes
}

// CleanLog removes the log file that contains possible crashes and the input
func CleanLog() {
	err := os.Remove("BerrFuzz-log.txt")
	if err != nil {
		fmt.Println("Failed to delete log file")
	}
}

// SetupLogger sets up local log file for fuzzing output //TODO: Possibly make logfilename input
func SetupLogger(cleanPtr *bool) {
	//! Check if log has to be cleaned
	if *cleanPtr {
		CleanLog()
	}

	logFile, err := os.OpenFile("BerrFuzz-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//defer logFile.Close()
	logFile.WriteString("\n")
	log.SetOutput(logFile)
	log.Printf("-=-==-===-====Start of New Fuzz Test")
}
