package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type balloon struct {
	top     string
	bottom  string
	borders [7]string
}

func newBalloon(width int) balloon {
	return balloon{
		top:     " " + strings.Repeat("_", width+2),
		bottom:  " " + strings.Repeat("-", width+2),
		borders: [7]string{"/", "\\", "\\", "/", "|", "<", ">"},
	}
}

func buildBalloon(lines []string, maxwidth int) string {
	b := newBalloon(maxwidth)
	var result []string
	result = append(result, b.top)

	if len(lines) == 1 {
		result = append(result, fmt.Sprintf("%s %s %s", b.borders[5], lines[0], b.borders[6]))
	} else {
		result = append(result, fmt.Sprintf("%s %s %s", b.borders[0], lines[0], b.borders[1]))
		for i := 1; i < len(lines)-1; i++ {
			result = append(result, fmt.Sprintf("%s %s %s", b.borders[4], lines[i], b.borders[4]))
		}
		result = append(result, fmt.Sprintf("%s %s %s", b.borders[2], lines[len(lines)-1], b.borders[3]))
	}

	result = append(result, b.bottom)
	return strings.Join(result, "\n")
}

func processLines(lines []string) ([]string, int) {
	// Replace tabs with spaces
	processed := make([]string, len(lines))
	maxWidth := 0

	for i, line := range lines {
		// Replace tabs with spaces
		line = strings.ReplaceAll(line, "\t", "    ")
		width := utf8.RuneCountInString(line)
		if width > maxWidth {
			maxWidth = width
		}
		processed[i] = line
	}

	// Normalize string lengths
	for i, line := range processed {
		processed[i] = line + strings.Repeat(" ", maxWidth-utf8.RuneCountInString(line))
	}

	return processed, maxWidth
}

var figures = map[string]string{

	"cow": `       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||`,

	"stegosaurus": `         \                      .       .
          \                    / ` + "`" + `.   .' "
           \           .---.  <    > <    >  .---.
            \          |    \  \ - ~ ~ - /  /    |
          _____           ..-~             ~-..-~
         |     |   \~~~\\.'                    ` + "`" + `./~~~/
        ---------   \__/                         \__/
       .'  O    \     /               /       \  "
      (_____,    ` + "`" + `._.'               |         }  \/~~~/
       ` + "`" + `----.          /       }     |        /    \__/
             ` + "`" + `-.      |       /      |       /      ` + "`" + `. ,~~|
                 ~-.__|      /_ - ~ ^|      /- _      ` + "`" + `..-'
                      |     /        |     /     ~-.     ` + "`" + `-. *  *  _
                      |_____|        |_____|         ~ - . *  *_>`,

	"kitty": `       \
        \
         \
          /l、
        （ﾟ､ ｡ 7
          l  ~ヽ
          じしf_,)/`,
}

func main() {
	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: fortune | gocowsay")
		return
	}

	figure := flag.String("f", "cow", "the figure name. Valid values are 'cow', 'stegosaurus', and 'kitty'")
	flag.Parse()

	// Read input lines
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return
	}

	// Process lines and build balloon
	processed, maxWidth := processLines(lines)
	balloon := buildBalloon(processed, maxWidth)

	fmt.Println(balloon)
	if art, ok := figures[*figure]; ok {
		fmt.Println(art)
	} else {
		fmt.Println("Unknown figure")
	}
	fmt.Println()
}
