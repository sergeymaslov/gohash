package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/martinlindhe/gohash"
)

var (
	encoding          = kingpin.Arg("encoding", "Output encoding.").String()
	listEncodings     = kingpin.Flag("list-encodings", "List available encodings.").Short('E').Bool()
	fileName          = kingpin.Arg("file", "Input file to read (optional).").String()
	encode            = kingpin.Flag("encode", "Encode input (default).").Short('e').Bool()
	decode            = kingpin.Flag("decode", "Decode input.").Short('d').Bool()
	outFileName       = kingpin.Flag("output", "Write output to file.").Short('o').String()
	noTrailingNewline = kingpin.Flag("no-newline", "Do not output the trailing newline.").Short('n').Bool()
	verbose           = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
)

func main() {

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *listEncodings {
		fmt.Println(gohash.AvailableEncodings())
		os.Exit(0)
	}

	encodings := strings.Split(*encoding, "+")

	if len(encodings) == 0 {
		fmt.Println("error: required argument 'encoding' not provided, try --help")
		os.Exit(1)
	}

	if *decode && *encode {
		fmt.Println("error: --decode and --encode don't mix")
		os.Exit(1)
	}

	r, err := gohash.ReadPipeOrFile(*fileName)
	if err != nil {
		log.Fatal("error:", err)
	}

	res, err := gohash.RecodeInput(encodings, r.Reader, *decode, *verbose)
	if err != nil {
		log.Fatal("error:", err)
	}

	if *outFileName != "" {
		f, err := os.Create(*outFileName)
		if err != nil {
			log.Fatal("error:", err)
		}
		defer f.Close()
		_, err = f.Write(res)
		if err != nil {
			log.Fatal("error:", err)
		}
		return
	}

	if *verbose {
		fmt.Println("final result:")
	}
	if *noTrailingNewline {
		fmt.Print(string(res))
	} else {
		fmt.Println(string(res))
	}
}
