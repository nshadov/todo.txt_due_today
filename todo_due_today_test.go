package main

import (
	"os"
	"regexp"
	"testing"
)

func TestGetTodayDate(t *testing.T) {
	today := getTodayDate()

	match, err := regexp.MatchString("20[0-9]{2}-[0-9]{2}-[0-9]{2}", today)
	if err != nil || match == false {
		t.Errorf("Invalid date format, got: %s, expected similar to: %s", today, "2018-01-31")
	}
}

func TestGetTodoFile(t *testing.T) {
	f := GetTodoFile()
	if f == nil {
		t.Errorf("Unable to open todo.txt file.")
	}
}

func TestFindTodayInFile(t *testing.T) {
	file, err := os.Open("./todo.txt")
	if err != nil {
		t.Errorf("Unable to open test todo file ./todo.txt")
	}
	defer file.Close()

	lines := FindDateInFile(file, "2018-01-31")
	if len(lines) != 1 {
		t.Errorf("Invalid number of matching events. Should be 1, but got: %s", lines)
	}
}
