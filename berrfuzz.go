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

// CheckTargetPath verifies or attempts to find the path of the target binary
func CheckTargetPath(targetBinName string) string {
	path, err := exec.LookPath(targetBinName)
	if err != nil {
		fmt.Printf("Target program '%s', not found in path. Provide full path!", targetBinName)
		os.Exit(-1)
	}
	return path
}

func main() {
	fmt.Println(string(ColorGreen), "-=BerrFuzz", string(ColorReset))

	cliPtr := flag.Bool("cli", false, "Use this flag for command line input as a payload")
	cleanPtr := flag.Bool("clean", false, "Use this flag to delete the log file and start with a fresh one")
	generatorPtr := flag.Bool("g", false, "Use this flag to put BerrFuzz in File Generation mode. This will generate files files")
	numTargerPtr := flag.Int("n", 0, `Use this flag to specify the number of iterations or number of files if used with -g flag. This flag
							 also requires flag 'i' and an input file`)
	mutationTypePtr := flag.String("m", "rand", `Use this flag to specify the mutation type you would like to use. The
										The default is rand`)
	exludeBytesPtr := flag.String("e", "", `Use this flag to exclude byte locations within the payload from being fuzzed 
									(i.e. header bytes). Examples of supported use:\n\tList: -e "0 1 2 7" -This will exclude bytes
									 0 1 2 and 7.\n\tRange: -e "0-4" - This will exclude bytes 0 through 4`)
	corpusFilenamePtr := flag.String("i", "", "Input Corpus file")
	targetBinName := flag.String("t", "", "Target program")
	corpusPayloadNamePtr := flag.String("c", "", "Name scheme for the output payload file")
	payloadSizePtr := flag.Float64("s", 0.0, `Given a number N < 1.0 , operations 
									will be done on N percent of the bytes.\n
									Given a number N >= 1.0, operations will be
									done on N many bytes\n`)
	replacePtr := flag.Bool("r", false, "Replace bytes with the payload rather than inserting the bytes")
	useOffsetPtr := flag.Bool("useoffset", false, "Tells BerrFuzz that you would like to insert the payload at a certain offset. Use with -o <int> flag for offset")
	offsetPtr := flag.Uint("o", 0, "Tells BerrFuzz what offset to insert the payload. This must be used in conjunction with the useoffset flag")

	flag.Parse()

	//TODO: Consider a mode selector
	//TODO: Logic for File Generator
	//! If *corpusFilenamePtr is empty string, then we will generate random bytes for files
	if *corpusFilenamePtr != "" && *numTargerPtr != 0 && *generatorPtr {
		FileGenerator(*corpusFilenamePtr, *numTargerPtr, *mutationTypePtr, *exludeBytesPtr, *payloadSizePtr, *useOffsetPtr, *offsetPtr, *replacePtr)
	}

	SetupLogger(cleanPtr)

	pathPrefix := OSCheck()

	// Corpus intake
	corpusFileBytes := make([]byte, 0)
	if *corpusFilenamePtr != "" {
		var corpusErr error
		corpusFileBytes, corpusErr = ReadFileBytes(*corpusFilenamePtr)
		if corpusErr != nil {
			fmt.Println("Error reading file: ", *corpusFilenamePtr)
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

	//! Harnessing Code

	*targetBinName = CheckTargetPath(*targetBinName) //TODO Move closer to the top

	//mutatedCorpusBytes := FileMutator(corpusFileBytes, payloadSize)

	//TODO: Error checking!
	if *corpusPayloadNamePtr == "" {
		fmt.Println("No corpus input provided")
	} else {
		// File Payload Using Corpus
		mutatedCorpusBytes := FileMutator(corpusFileBytes, payloadSize)
		OutputCorpus(*corpusPayloadNamePtr, mutatedCorpusBytes)
		*corpusPayloadNamePtr = pathPrefix + *corpusPayloadNamePtr
	}

	payload := ""

	if *cliPtr {
		payload = string(RandomByteGenerator(payloadSize))
	} else {
		payload = *corpusPayloadNamePtr
	}

	payload = "<< " + payload

	cmd := exec.Command(*targetBinName, payload)

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
}
