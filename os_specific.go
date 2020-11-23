package main

import (
	"fmt"
	"runtime"
)

// OSCheck checks the operating system for OS specific functionality
func OSCheck() (string, string) {
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
	return pathPrefix, runtime.GOOS
}
