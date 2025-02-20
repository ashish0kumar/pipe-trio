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
			if lines[i] == "" {
				result = append(result, fmt.Sprintf("%s %s %s",
					b.borders[4],
					strings.Repeat(" ", maxwidth),
					b.borders[4]))
			} else {
				result = append(result, fmt.Sprintf("%s %s %s", b.borders[4], lines[i], b.borders[4]))
			}
		}

		result = append(result, fmt.Sprintf("%s %s %s", b.borders[2], lines[len(lines)-1], b.borders[3]))
	}

	result = append(result, b.bottom)
	return strings.Join(result, "\n")
}

func processLines(lines []string) ([]string, int) {
	processed := make([]string, len(lines))
	maxWidth := 0
	for i, line := range lines {
		line = strings.ReplaceAll(line, "\t", "    ")
		width := utf8.RuneCountInString(line)
		if width > maxWidth {
			maxWidth = width
		}
		processed[i] = line
	}

	for i, line := range processed {
		if line != "" {
			processed[i] = line + strings.Repeat(" ", maxWidth-utf8.RuneCountInString(line))
		}
	}
	return processed, maxWidth
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine string
	currentLineWidth := 0

	for _, word := range words {
		wordWidth := utf8.RuneCountInString(word)

		if currentLineWidth > 0 && currentLineWidth+wordWidth+1 > width {
			lines = append(lines, strings.TrimSpace(currentLine))
			currentLine = word
			currentLineWidth = wordWidth
		} else {
			if currentLineWidth > 0 {
				currentLine += " "
				currentLineWidth++
			}
			currentLine += word
			currentLineWidth += wordWidth
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return lines
}

// handles paragraphs and preserves some formatting
func enhancedWrapText(text string, width int) []string {
	paragraphs := strings.Split(text, "\n\n")
	var result []string

	for i, paragraph := range paragraphs {
		lines := strings.Split(paragraph, "\n")
		var paragraphText string

		for _, line := range lines {
			if len(paragraphText) > 0 {
				paragraphText += " "
			}
			paragraphText += strings.TrimSpace(line)
		}

		wrappedParagraph := wrapText(paragraphText, width)
		result = append(result, wrappedParagraph...)

		if i < len(paragraphs)-1 {
			result = append(result, "")
		}
	}

	return result
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
	wrapWidth := flag.Int("w", 40, "the maximum width for text wrapping")
	flag.Parse()

	// Read all input
	reader := bufio.NewReader(os.Stdin)
	var inputText strings.Builder
	for {
		line, err := reader.ReadString('\n')
		inputText.WriteString(line)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return
		}
		if err == io.EOF {
			break
		}
	}

	// Wrap the text
	wrappedLines := enhancedWrapText(inputText.String(), *wrapWidth)

	// Process lines and build balloon
	processed, maxWidth := processLines(wrappedLines)
	balloon := buildBalloon(processed, maxWidth)

	fmt.Println(balloon)
	if art, ok := figures[*figure]; ok {
		fmt.Println(art)
	} else {
		fmt.Println("Unknown figure")
	}
	fmt.Println()
}
