// BerrFuzz v.0001 ;)
//        by
// Manuel Berrueta

package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"os"
	"os/exec"
	"runtime"
)

// Color for color for your visual pleasure :)
type Color string

// Some color for fun!
const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)

// DEBUG  global var for debugging
const DEBUG = true

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

// SetupLogger sets up local log file for fuzzing output //TODO: Possibly make logfilename input
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
	randByte := mrand.Intn(0xFFFFFFFF)

	for i := 0; i < numOfByteToFlip; i++ {
		fileBytes[mrand.Intn(len(fileBytes))-1] = byte(randByte)
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

func main() {
	fmt.Println(string(ColorGreen), "-=BerrFuzz", string(ColorReset))

	cleanPtr := flag.Bool("clean", false, "Delete log file")
	corpusName := ""
	flag.StringVar(&corpusName, "i", "", "Input Corpus file")
	targetName := ""
	flag.StringVar(&targetName, "t", "", "Target program")
	corpusPayloadName := ""
	flag.StringVar(&corpusPayloadName, "c", "", "Name for the output payload file")

	payloadSizePtr := flag.Float64("s", 0.0, `Given a number N < 1.0 , operations 
									will be done on N percent of the bytes.\n
									Given a number N >= 1.0, operations will be
									done on N many bytes\n`)

	flag.Parse()

	if *cleanPtr {
		CleanLog()
	}

	//Testing Flags + Parameters
	//fmt.Println("flag.Args()[0]:", flag.Args()[0])

	SetupLogger()

	// OS Check run for running compatible commands
	// WIP
	//OS := ""
	pathPrefix := ""
	switch runtime.GOOS {
	case "windows":
		//ver, err := syscall.GetVersion()
		//if err != nil {
		//	panic(err)
		//}
		//Major := int(ver & 0xFF)
		//Minor := int(ver >> 8 & 0xFF)
		//Build := int(ver >> 16)
		//OS = "windows"
		fmt.Printf("Running on Windows %d Build: %d | Arch: %s | CPU(s): %d\n", 0, 0, runtime.GOARCH, runtime.NumCPU()) //!WIP
		pathPrefix = ".\\"
	case "linux":
		fmt.Printf("Running on Linux '%s' | Ver: %s | Arch: %s | CPU(s): %d\n", "Ubuntu/Fedora/Whatever", "verTBD", runtime.GOARCH, runtime.NumCPU())
		//OS = "linux"
		pathPrefix = "./"
	case "darwin":
		fmt.Printf("Running on Mac OS '%s' | Ver: %s | Arch: %s | CPU(s): %d\n", "?", "verTBD", runtime.GOARCH, runtime.NumCPU())
		//OS = "darwin"
		pathPrefix = "./"
	}

	// Corpus intake
	corpusFileBytes := make([]byte, 0)
	if corpusName != "" {
		var corpusErr error
		corpusFileBytes, corpusErr = ReadFileBytes(corpusName)
		if corpusErr != nil {
			fmt.Println("Error reading file: ", corpusName)
			os.Exit(-1)
		}
	}

	payloadSize := 0
	if *payloadSizePtr != 0.0 {
		/* If payloadSizePtr is < 1.0 we should have an input file and we will
		   calculate the size as a percentage of that file size */
		if *payloadSizePtr < 1.0 && *payloadSizePtr > 0.0 {
			//* Check if a filename is passed
			if len(corpusFileBytes) == 0 {
				fmt.Println("Error: Passed percentage but no input file!")
				os.Exit(-1)
			} else {
				payloadSize = int(float64(len(corpusFileBytes)) * *payloadSizePtr)
			}
		} else if *payloadSizePtr >= 1.0 {
			payloadSize = int(*payloadSizePtr)
		} else {
			fmt.Println("Payload size must be a positive number!")
			os.Exit(-1)
		}
	} else {
		//payloadSize = mrand.Intn(0xFFFFFFFF)
		payloadSize = mrand.Intn(0x1fffffe8)
	}

	//Testing payload generation
	if DEBUG {
		fmt.Println("Randomly generated payload size: ", payloadSize)
		payloadSize = 2000
		fmt.Println("Debugging size: ", payloadSize)
		payload := string(RandomByteGenerator(payloadSize))
		payload2 := string(RandomBitFlip(RandomByteGenerator(payloadSize)))
		fmt.Println("\nPayload_2: ", payload2)
		fmt.Println("\nPayload: ", payload)
	}

	//Find target in path
	path, err := exec.LookPath(targetName)
	if err != nil {
		fmt.Printf("Target program '%s', not found in path. Provide full path!", targetName)
		os.Exit(-1)
	}
	targetName = path

	mutatedCorpusBytes := FileMutator(corpusFileBytes, payloadSize)

	//TODO: Error checking!
	if corpusPayloadName == "" {
		fmt.Println("No corpus input provided")
	}
	// File Payload Using Corpus
	OutputCorpus(corpusPayloadName, mutatedCorpusBytes)

	//cmd := exec.Command(*targetPtr)
	//cmd := exec.Command(targetName, string(payload))

	corpusPayloadName = pathPrefix + corpusPayloadName

	cmd := exec.Command(targetName, corpusPayloadName)

	//Output redirections stdout and stderr
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	log.Printf("Running command")

	output, err := cmd.CombinedOutput()
	if err != nil {

		log.Printf("-=-=####### Possible Crash:")
		log.Printf("-=-=####### \tErr: %s", err.Error())
		//log.Printf("-=-=#######\t\tPayload @ Possible Crash:\n \"%s\"", payload)
		//TODO: Create a dir Possible_Crashes, copy payload/payload file to it
	}

	fmt.Printf("Output of program: %s\n", string(output))
	log.Println("Done running command")

	fmt.Println("Testing RandomFileGeneration")
	RandomByteFileGenerator(1024, "randomFile")

	//TODO: Being able to choose certain character sets
	//TODO: Integrate known bad strings
	//TODO: Possibly integrate search for known bad functions
	//TODO: Payload can be an optional input
}
