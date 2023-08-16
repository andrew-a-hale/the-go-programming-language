package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unitconv/tempconv"
	"unitconv/weightconv"
)

var (
	c                 = flag.String("c", "", "type of conversion")
	tempconversions   = []string{"ctof", "ctok", "ftoc", "ftok", "ktoc", "ktof"}
	weightconversions = []string{"lbtokg", "kgtolb"}
)

func main() {
	flag.Parse()

	in, err := strconv.ParseFloat(flag.Arg(0), 64)
	if err != nil {
		log.Print("enter a number to convert: ")
		buf := bufio.NewReader(os.Stdin)
		in_text, _ := buf.ReadBytes('\n')
		if match, err := regexp.Match("[0-9]+", in_text); !match || err != nil {
			log.Fatal("bad input")
		}
		in = float64(in_text[0])
	}

	if cIsA(*c, tempconversions) {
		tempConverter(*c, in)
	} else if cIsA(*c, weightconversions) {
		weightConverter(*c, in)
	} else {
		log.Fatal("bad c arg")
	}
}

func cIsA(c string, ctypes []string) bool {
	for _, v := range ctypes {
		if c == v {
			return true
		}
	}
	return false
}

func tempConverter(c string, in float64) {
	switch strings.ToLower(c) {
	case "ctof":
		input := tempconv.Celsius(in)
		formatAndPrint(input.String(), tempconv.CToF(input).String())
	case "ctok":
		input := tempconv.Celsius(in)
		formatAndPrint(input.String(), tempconv.CToK(input).String())
	case "ftoc":
		input := tempconv.Fahrenheit(in)
		formatAndPrint(input.String(), tempconv.FToC(input).String())
	case "ftok":
		input := tempconv.Fahrenheit(in)
		formatAndPrint(input.String(), tempconv.FToK(input).String())
	case "ktoc":
		input := tempconv.Kelvin(in)
		formatAndPrint(input.String(), tempconv.KToC(input).String())
	case "ktof":
		input := tempconv.Kelvin(in)
		formatAndPrint(input.String(), tempconv.KToF(input).String())
	case "":
		log.Fatal("missing conversion type")
	}
}

func weightConverter(c string, in float64) {
	switch strings.ToLower(c) {
	case "lbtokg":
		input := weightconv.Pound(in)
		formatAndPrint(input.String(), weightconv.LbToKg(input).String())
	case "kgtolb":
		input := weightconv.Kilogram(in)
		formatAndPrint(input.String(), weightconv.KgToLb(input).String())
	}
}

func formatAndPrint(s ...string) {
	fmt.Println(strings.Join(s, " => "))
}
