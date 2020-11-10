package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func OffsetByteInsert(inBytes []byte, payload []byte, offset uint, replace bool) {

}

func KeyByteInsert(inBytes []byte, payload []byte, keyTarget string, replace bool) {
	//! Search for matching keyTarget
	//index := bytes.Index(inBytes, []byte(keyTarget))
	if replace {
		bytes.Replace(inBytes, []byte(keyTarget), payload, 1) //n = number of replacements
	} else {
		index := bytes.Index(inBytes, []byte(keyTarget))
		//TODO: Either split the InBytes at the index in to two and append payload to first part and second part after payload
		//TODO:  or find alternative way to insert at that location of bytes.
	}

}

// FileGenerator creates files depending on the input parameters
func FileGenerator(inputCorpusFilename string, numOfFiles int, mutationType string, excludedBytes string, payloadSizeFl float64, useOffset bool, offset uint, replace bool) {
	//TODO: Logic for File Generator

	//TODO: Extensions support, to append an extension to the outfilename

	nameSchem := "fuzzInputFile"

	//! We do not exclude any bytes from fuzzing
	if excludedBytes != "" {
		//TODO: we keep this bytes in a temp structure or use some logic to exclude them
		//* Examples: Given a range from 0-20 (possibly header bytes) keep temp bytes and replace after this bytes after mutation
	}

	corpusFileBytes := make([]byte, 0)
	originalBytes := make([]byte, 0)

	//! If a input corpus filename is given, we open and read in the bytes
	if inputCorpusFilename != "" {
		corpusFileBytes = CorpusIntake(inputCorpusFilename)
		originalBytes = corpusFileBytes //! Save this for manipulation of original bytes
	}

	//! Calculate payload size based on the payloadSizeFl parameter
	payloadSize := CalculatePayloadSize(payloadSizeFl, corpusFileBytes)

	//! Create files
	for i := 0; i < numOfFiles; i++ {
		if len(originalBytes) != 0 {
			corpusFileBytes = originalBytes
			// Byte mutation

			//TODO: Select mutation type
			var mutatedBytes []byte
			switch mutationType {
			case "":
				fallthrough
			case "rand":
				if useOffset == false {
					mutatedBytes = FileMutator(corpusFileBytes, payloadSize)
				} else {
					//todo: create random bytes
					if replace {
						//todo: here we overwrite payloadsize original bytes
					} else {
						//todo: here we concatenate the payload with the original bytes
					}
				}
			case "input":
				//todo: this is the case where a payload is passed
			case "keyTarget":
				//todo: this is the case where we put a byte pattern in the file as our target
			}

			//TODO: Here we put together file exluded bytes

			OutputCorpus(nameSchem+strconv.Itoa(i), mutatedBytes)

		} else {

			RandomByteFileGenerator(payloadSize, nameSchem+strconv.Itoa(i))
		}
	}

	fmt.Println("Finished Creating Files!")
	os.Exit(0)
}
