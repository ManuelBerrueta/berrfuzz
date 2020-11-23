package main

import (
	"fmt"
	"math/bits"
	mrand "math/rand"
)

// BitFlipper will flip the byteIndex bit in argument inByte, where index 0 is the MSB
func BitFlipper(byteIndex int, inByte byte) byte {
	if byteIndex < 0 || byteIndex > 7 {
		fmt.Println("Error: BitFlipper() requires a byteIndex in range of 0-7")
	}

	for byteIndex < 0 || byteIndex > 7 {
		fmt.Printf("Enter byteIndex in range of 0-7: ")
		fmt.Scan(&byteIndex)
	}

	// Could also use inByte |= (1 << position) with a bit different logic...
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
	case 7: //! LSB
		byteWithFlippedBit = inByte ^ 1
	}

	return byteWithFlippedBit
}

// LeftByteShift shifts inByte by the amount of the  argument shiftBy
func LeftByteShift(shiftBy int, inByte byte) byte {
	shiftedLeftByte := inByte << shiftBy
	return shiftedLeftByte
}

// RightByteShift shifts inByte by the amount of the  argument shiftBy
func RightByteShift(shiftBy int, inByte byte) byte {
	shiftedRightByte := inByte >> shiftBy
	return shiftedRightByte
}

// RandomBitFlip will flip just 1 bit in a random Byte
func RandomBitFlip(inData []byte) []byte {
	randByteIndex := mrand.Intn(len(inData) - 1)
	tempRandByte := inData[randByteIndex]
	tempRandByte = BitFlipper(mrand.Intn(7), tempRandByte)
	inData[randByteIndex] = tempRandByte

	return inData
}

// RandomBitFlips flips a numBitsToFlip number of bits randomly
func RandomBitFlips(inData []byte, numBitsToFlip int) []byte {
	for i := 0; i < numBitsToFlip; i++ {
		randByteIndex := mrand.Intn(len(inData) - 1)
		tempRandByte := inData[randByteIndex]
		tempRandByte = BitFlipper(mrand.Intn(7), tempRandByte)
		inData[randByteIndex] = tempRandByte
	}
	return inData
}

// ByteFlipper flips all the bytes in the passed argument byte
func ByteFlipper(byteToFlip byte) byte {
	flippedBytes := ^byteToFlip
	return flippedBytes
}

// ReverseByte reverses all the bits on the passed byteToReverse byte
func ReverseByte(byteToReverse byte) byte {
	reversedByte := bits.Reverse(uint(byteToReverse))
	return byte(reversedByte)
}

//* Strings ops

func strExists(inStr []string, target string) (bool, int) {
	for i, char := range inStr {
		if char == target {
			return true, i
		}
	}
	return false, -1
}
