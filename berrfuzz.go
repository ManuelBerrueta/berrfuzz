package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
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

// ReadFileBytes reads a file and returns bytes
func ReadFileBytes(fileName string) ([]byte, error) {
	inFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

// BitFlipper will flip the byteIndex bit in argument inByte, where index 0 is the MSB
func BitFlipper(byteIndex int, inByte byte) byte {
	if byteIndex < 0 || byteIndex > 7 {
		fmt.Println("Error: BitFlipper() requires a byteIndex in range of 0-7")
	}

	for byteIndex < 0 || byteIndex > 7 {
		fmt.Printf("Enter byteIndex in range of 0-7: ")
		fmt.Scan(&byteIndex)
	}

	var byteWithFlippedBit byte
	switch byteIndex {
	case 0: //! Flip MSB
		byteWithFlippedBit = inByte ^ 128
	case 1:
		byteWithFlippedBit = inByte ^ 64
	case 2:
		byteWithFlippedBit = inByte ^ 32
	case 3:
		byteWithFlippedBit = inByte ^ 16
	case 4:
		byteWithFlippedBit = inByte ^ 8
	case 5:
		byteWithFlippedBit = inByte ^ 4
	case 6:
		byteWithFlippedBit = inByte ^ 2
	case 7:
		byteWithFlippedBit = inByte ^ 1
	}
	return byteWithFlippedBit
}

// LeftByteShift shifts inByte by the amount of the  argument shiftBy
func LeftByteShift(shiftBy int, inByte byte) byte {
	shiftedLeftByte := inByte << shiftBy
	return shiftedLeftByte
}

// RandomBitFlip will flip just 1 bit in a random Byte
func RandomBitFlip(inData []byte, numBitsToFlip int) []byte {
	randByteIndex := mrand.Intn(len(inData) - 1)
	tempRandByte := inData[randByteIndex]

	// This is a byteshift
	inData[randByteIndex] = tempRandByte << 1

	//Select random bit
	//rndBitIndex := mrand.Intn(7//)
	//tempBit := tempRandByte[ran//dBitIndex]

	//flippedBit := ^randBit
	//inData[randByte][randBit] = flippedBit
	//TODO: Need to work on bit selection

	return inData
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

// Generator()
// --- This will generate a random file, which will then be output to be used by the fuzzer
//Mutator() This will mutate the existing file
// --- Should eventually develop into something that allows us to select what to mutate
// FileMutator is the start of a more complex function that will use different
//             techniques to flip bits/bytes
func FileMutator(fileBytes []byte) {
	numOfByteToFlip := 10

	//C
	randByte := mrand.Intn(0xFFFFFFFF)

	for i := 0; i < numOfByteToFlip; i++ {
		fileBytes[mrand.Intn(len(fileBytes))-1] = byte(randByte)
	}
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
	flag.Parse()

	if *cleanPtr {
		CleanLog()
	}

	//Testing Flags + Parameters
	fmt.Println("flag.Args:", flag.Args()[0])
	fmt.Println("Args:", os.Args)

	//Add options to delete local log file
	SetupLogger()

	// OS Check run for running compatible commands
	switch runtime.GOOS {
	case "windows":
		//ver, err := syscall.GetVersion()
		//if err != nil {
		//	panic(err)
		//}
		//Major := int(ver & 0xFF)
		//Minor := int(ver >> 8 & 0xFF)
		//Build := int(ver >> 16)
		fmt.Printf("Running on Windows %d Build: %d | Arch: %s | CPU(s): %d\n", 0, 0, runtime.GOARCH, runtime.NumCPU()) //!WIP
	case "linux":
		fmt.Printf("Running on Linux '%s' | Ver: %s | Arch: %s | CPU(s): %d\n", "Ubuntu/Fedora/Whatever", "verTBD", runtime.GOARCH, runtime.NumCPU())
	case "darwin":
		fmt.Printf("Running on Mac OS '%s' | Ver: %s | Arch: %s | CPU(s): %d\n", "?", "verTBD", runtime.GOARCH, runtime.NumCPU())
	}

	fileBytes, err := ReadFileBytes("test.txt")
	if err != nil {
		fmt.Println("Error reading files")
	}
	fmt.Println(string(fileBytes))

	totalNum := 2000

	//! Could make an additional arg to be passed depends on what is beingtotalNum := 2000
	payload := string(RandomByteGenerator(totalNum))
	println("Payload: ", payload)

	//! This will be a command line input
	targetProgram := "powershell.exe"

	cmd := exec.Command(targetProgram)
	//cmd := exec.Command(targetProgram, string(payload))
	//cmd := exec.Command(targetProgram)

	//Output redirections stdout and stderr
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	log.Printf("Running command")

	output, err := cmd.CombinedOutput()
	if err != nil {

		log.Printf("-=-=####### Possible Crash:")
		log.Printf("-=-=####### \tErr: %s", err.Error())
		log.Printf("-=-=#######\t\tPayload @ Possible Crash:\n \"%s\"", payload)
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
