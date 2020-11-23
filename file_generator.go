package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// OffsetByteInsert will insert the payload at the given byte offset
func OffsetByteInsert(inBytes []byte, payload []byte, offset uint) []byte {
	appendedBytes := append(inBytes[:offset], payload[:]...)
	//! Append the rest of the inBytes from the index to the back of the payload
	appendedBytes = append(appendedBytes, inBytes[offset+1:]...)
	return appendedBytes
}

// KeyByteInsert will insert the payload at the given keyTarget bytes location
func KeyByteInsert(inBytes []byte, payload []byte, keyTarget string, replace bool) []byte {
	if replace {
		bytes.Replace(inBytes, []byte(keyTarget), payload, 1) //n = number of replacements
	} else {
		index := bytes.Index(inBytes, []byte(keyTarget))
		//TODO: Either split the InBytes at the index in to two and append payload to first part and second part after payload
		//TODO:  or find alternative way to insert at that location of bytes.
		//inBytes = append(payload,)
		//! Append the payload to the inBytes at desired index
		appendedBytes := append(inBytes[:index-1], payload[:]...)
		//! Append the rest of the inBytes from the index to the back of the payload
		appendedBytes = append(appendedBytes, inBytes[(index+len(keyTarget)):]...)
		inBytes = appendedBytes
	}
	return inBytes
}

// ByteRangeCutter removes a certain number (numOfBytesToCut) of bytes from inBytes
func ByteRangeCutter(inBytes []byte, offset uint, numOfBytesToCut uint) []byte {
	cutBytes := inBytes[:offset-1]
	cutBytes = append(cutBytes, inBytes[(offset+numOfBytesToCut):]...)
	return cutBytes
}

// ByteTargetTrimmer trims/removes a certain number (numOfInstancestoCut) of the targetBytes from inBytes
func ByteTargetTrimmer(inBytes []byte, targetBytes []byte, numOfInstancesToCut int) []byte {
	for i := 0; i < numOfInstancesToCut; i++ {
		inBytes = bytes.Trim(inBytes, string(targetBytes))
	}
	return inBytes
}

// ByteTargetClipper removes all instances of targetBytes from inBytes
func ByteTargetClipper(inBytes []byte, targetBytes []byte) []byte {
	inBytes = bytes.TrimRight(inBytes, string(targetBytes))
	return inBytes
}

// FileGenerator creates files with a variety of mutations depending on the input parameters
func FileGenerator(settings Settings) {

	//TODO: Extensions support, to append an extension to the outfilename

	nameSchem := "fuzzInputFile"

	//! We do not exclude any bytes from fuzzing
	if settings.ExcludedBytes != "" {
		//TODO: we keep this bytes in a temp structure or use some logic to exclude them
		//* Examples: Given a range from 0-20 (possibly header bytes) keep temp bytes and replace after this bytes after mutation

		//* Check if ExludeBytes includes a range or individual bytes

		//TODO: work on excluding bytes
		if strings.Contains(settings.ExcludedBytes, "-") { //We are dealing with are range of bytes
			//delimOffset := strings.Index(settings.ExcludedBytes, "-")
			excludeRangeStr := strings.Split(settings.ExcludedBytes, "-")

			var excludeRange []int
			for i, num := range excludeRangeStr {
				excludeRange[i], _ = strconv.Atoi(num)
			}
			//? excludeRange includes the right and left numbers for the range
		}
	}

	corpusFileBytes := make([]byte, 0)
	originalBytes := make([]byte, 0)

	//! If a input corpus filename is given, we open and read in the bytes

	if settings.InputCorpusFilename != "" {
		corpusFileBytes = CorpusIntake(settings.InputCorpusFilename)
		originalBytes = corpusFileBytes //! Save this for manipulation of original bytes
	}

	//! Calculate payload size based on the payloadSizeFl parameter
	payloadSize := CalculatePayloadSize(settings.PayloadSizeFl, corpusFileBytes)

	//! Create files
	for i := 0; i < settings.NumTarget; i++ {
		if len(originalBytes) != 0 {
			corpusFileBytes = originalBytes
			// Byte mutation

			//TODO: Select mutation type
			var mutatedBytes []byte
			switch settings.MutationType {
			case "":
				fallthrough
			case "rand":
				if settings.UseOffset == false {
					mutatedBytes = FileMutator(corpusFileBytes, payloadSize)
				} else {
					//todo: create random bytes
					if settings.Replace {
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
