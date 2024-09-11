package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	// Define the flags for the input script and output file
	scriptPath := flag.String("input-script", "", "Path to the input script")
	outputPath := flag.String("output-to", "", "Path to the output ELF file")
	flag.Parse()

	// Check if both flags are provided
	if *scriptPath == "" || *outputPath == "" {
		fmt.Println("Error: --input-script and --output-to flags are required")
		flag.Usage()
		os.Exit(1)
	}

	// Read the input script
	scriptContent, err := ioutil.ReadFile(*scriptPath)
	if err != nil {
		fmt.Printf("Error reading script file: %v\n", err)
		os.Exit(1)
	}

	// Define the ELF header bytes
	elfHeader := []byte{
		0x7f, 0x45, 0x4c, 0x46, // ELF magic number
		0x02, 0x01, 0x01, 0x00, // 64-bit, little-endian, ELF version 1
	}

	// Create the output file
	file, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write the ELF header to the output file
	_, err = file.Write(elfHeader)
	if err != nil {
		fmt.Printf("Error writing ELF header to file: %v\n", err)
		os.Exit(1)
	}

	// Define the regular expression for matching a shebang at the beginning of the script
	shebangPattern := regexp.MustCompile(`^#!\s*(/.*?)(?:\s|$)`)
	shebangMatch := shebangPattern.Find(scriptContent)

	// Check if the script starts with a shebang
	if shebangMatch == nil {
		fmt.Println("Error: The script must have a shebang as the first line.")
		os.Exit(1)
	}

	// Write the shebang to the output file
	_, err = file.Write(shebangMatch)
	if err != nil {
		fmt.Printf("Error writing shebang to file: %v\n", err)
		os.Exit(1)
	}

	// Add the printf statement right after the shebang
	_, err = file.WriteString("\n" + `printf "\x1B[1F\x1B[2K"`)
	if err != nil {
		fmt.Printf("Error writing printf command to file: %v\n", err)
		os.Exit(1)
	}

	// Write the rest of the script content (excluding the shebang part)
	_, err = file.Write(scriptContent[len(shebangMatch):])
	if err != nil {
		fmt.Printf("Error writing script content to file: %v\n", err)
		os.Exit(1)
	}

	// Make the output file executable
	err = os.Chmod(*outputPath, 0755)
	if err != nil {
		fmt.Printf("Error setting file permissions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created ELF file with embedded script: %s\n", *outputPath)
}
