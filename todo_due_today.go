package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

var rexDate, _ = regexp.Compile("due:[0-9]{4}-[0-9]{2}-[0-9]{2}")

func getTodayDate() string {
	t := time.Now()
	return fmt.Sprintf("%4d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

// GetTodoFile return file descriptor to todo.txt file
func GetTodoFile() *os.File {
	cmd := "source ~/.todo/config && echo ${TODO_FILE}"

	subprocess := exec.Command("/bin/sh")
	stdin, err := subprocess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(stdin, cmd)
	stdin.Close()
	subprocess.Wait()
	out, err := subprocess.Output()
	if err != nil {
		log.Fatal(out, err)
	}

	file, err := os.Open(strings.Trim(string(out), "\n"))
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// FindDateInFile searches lines that contains date older than specified
func FindDateInFile(file *os.File, date string) []string {
	scanner := bufio.NewScanner(file)
	output := []string{}

	defer file.Close()
	timeNow, _ := time.Parse("2006-01-02", date)

	for scanner.Scan() {
		line := scanner.Text()
		match := rexDate.FindString(line)
		if match != "" {
			timeEvent, _ := time.Parse("due:2006-01-02", match)
			if timeEvent.Before(timeNow) || timeEvent.Equal(timeNow) {
				output = append(output, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return output
}

// FilterContext checks if specified line contains context
func FilterContext(line, context string) *string {
	if strings.Contains(line, context) {
		return &line
	}
	return nil
}

// IsFinished answers if task already has been done
func IsFinished(line string) bool {
	if strings.HasPrefix(line, "x ") {
		return true
	}
	return false
}

// PrintColorText displays color lines
func PrintColorText(line string) {
	if strings.HasPrefix(line, "(A)") {
		ct.Foreground(ct.Yellow, false)
	} else if strings.HasPrefix(line, "(B)") {
		ct.Foreground(ct.Green, false)
	} else if strings.HasPrefix(line, "(C)") {
		ct.Foreground(ct.Blue, false)
	}
	fmt.Println(line)
	ct.ResetColor()
}

func main() {
	var context string

	if len(os.Args) > 1 {
		context = os.Args[1]
	} else {
		context = "@work"
	}

	todo := FindDateInFile(GetTodoFile(), getTodayDate())
	sort.Strings(todo)
	for _, line := range todo {
		if FilterContext(line, context) != nil && IsFinished(line) == false {
			PrintColorText(line)
		}
	}
}
