package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

func rgb(i int) (int, int, int) {
	f := 0.1
	phase1 := 0.0
	phase2 := 2 * math.Pi / 3
	phase3 := 4 * math.Pi / 3

	return int(math.Sin(f*float64(i)+phase1)*127 + 128),
		int(math.Sin(f*float64(i)+phase2)*127 + 128),
		int(math.Sin(f*float64(i)+phase3)*127 + 128)
}

func main() {
	if isTerminal() {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: fortune | gololcat")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	colorIndex := 0

	for {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		r, g, b := rgb(colorIndex)
		fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, char)
		colorIndex++
	}
}

func isTerminal() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice != 0
}
