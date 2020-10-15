// BerrFuzz v.0001 ;)
//        by
// Manuel Berrueta

package main

import (
	"flag"
	"fmt"
	"log"
	mrand "math/rand"
	"os"
	"os/exec"
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
const DEBUG = false

// checkTargetPath verifies or attempts to find the path of the target binary
func checkTargetPath(targetName string) string {
	path, err := exec.LookPath(targetName)
	if err != nil {
		fmt.Printf("Target program '%s', not found in path. Provide full path!", targetName)
		os.Exit(-1)
	}
	return path
}

func main() {
	fmt.Println(string(ColorGreen), "-=BerrFuzz", string(ColorReset))

	cliPtr := flag.Bool("cli", false, "Use this flag for command line input only")
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

	SetupLogger(cleanPtr)

	pathPrefix := OSCheck()

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

	//TODO: This is only considering the case where we have an input file, what if we want to shove input to a program

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

	targetName = checkTargetPath(targetName) //TODO Move closer to the top

	//mutatedCorpusBytes := FileMutator(corpusFileBytes, payloadSize)

	//TODO: Error checking!
	if corpusPayloadName == "" {
		fmt.Println("No corpus input provided")
	} else {
		// File Payload Using Corpus
		mutatedCorpusBytes := FileMutator(corpusFileBytes, payloadSize)
		OutputCorpus(corpusPayloadName, mutatedCorpusBytes)
		corpusPayloadName = pathPrefix + corpusPayloadName
	}

	payload := ""

	if *cliPtr {
		payload = string(RandomByteGenerator(payloadSize))
	} else {
		payload = corpusPayloadName
	}

	payload = "<< " + payload

	cmd := exec.Command(targetName, payload)

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

	//TODO: Being able to choose certain character sets
	//TODO: Integrate known bad strings
	//TODO: Possibly integrate search for known bad functions
	//TODO: Payload can be an optional input
}
