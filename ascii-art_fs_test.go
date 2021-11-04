package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

/*	Each key of the map testCases contains the name of the file containing the
	expected out put for each test case, the value for each key is a slice of
	strings, the first element contains the string to be given as an argument
	at run time, the second will contain the string equivalent of expected
	output	*/
var testCases = map[int][]string{
	1:  {"banana", "standard", "asd", ""},
	2:  {"hello", "standard", ""},
	3:  {"hello world", "shadow", ""},
	4:  {"nice 2 meet you", "thinkertoy", ""},
	5:  {"you & me", "standard", ""},
	6:  {"123", "shadow", ""},
	7:  {"/(\")", "thinkertoy", ""},
	8:  {"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "shadow", ""},
	9:  {"\"#$%&/()*+,-./", "thinkertoy", ""},
	10: {"It's Working", "thinkertoy", ""},
}

/*	This test file tests the ascii-art project against the first 9 test cases on
audit page	*/
func TestAsciiArtFS(t *testing.T) {
	getTestCases()

	// Test the program with incorrect amount of args
	output, err := exec.Command("go", "run", ".", testCases[1][0], testCases[1][1], testCases[1][2]).Output()
	if err != nil {
		panic(err)
	}
	if string(output) != testCases[1][3] {
		t.Errorf("\nTest fails when given the arguments:\n\t\"%s\" \"%s\" \"%s\","+
			"\nexpected:\n%s\ngot:\n%s\n\n",
			testCases[1][0], testCases[1][1], testCases[1][2], testCases[1][3], string(output))
	}

	/*	Iterate through each test case and starting a goroutine for each, this
		is done so instead of waiting for the previous test to complete they can
		all be checked simulaneously	*/
	var wg sync.WaitGroup
	for i := 2; i <= len(testCases); i++ {
		wg.Add(1)
		go func(current []string, w *sync.WaitGroup, ti *testing.T) {
			defer w.Done()
			result := getResult(current)
			/*	Fails the project if the test cases expected output doesn't match
				the actual output	*/
			if string(result) != current[2] {
				ti.Errorf("\nTest fails when given the test case:\n\t\"%s\" \"%s\","+
					"\nexpected:\n%s\ngot:\n%s\n\n",
					current[0], current[1], current[2], string(result))
			}
		}(testCases[i], &wg, t)
	}
	wg.Wait()
}

/*	This function imitates the running of "go run . string", which it then pipes
	into a second function "cat -e" to immitate and then returns the result	*/
func getResult(testCase []string) string {
	first := exec.Command("go", "run", ".", testCase[0], testCase[1])
	second := exec.Command("cat", "-e")
	reader, writer := io.Pipe()
	first.Stdout = writer
	second.Stdin = reader
	var buffer bytes.Buffer
	second.Stdout = &buffer
	first.Start()
	second.Start()
	first.Wait()
	writer.Close()
	second.Wait()
	return buffer.String()
}

/*	This function reads each of the test cases expected output from the "testcases.txt"
	file and adds them to the corresponding test cases slice in the testCases map	*/
func getTestCases() {
	file, err := os.Open("test-cases.txt")
	if err != nil {
		panic(err)
	}

	stats, _ := file.Stat()
	contents := make([]byte, stats.Size())
	file.Read(contents)
	lines := strings.Split(string(contents), "\n")

	start := 0
	number := 0
	for i, line := range lines {
		if i == len(lines)-1 {
			testCases[number][len(testCases[number])-1] = strings.Join(lines[start:], "\n") + "\n"
			break
		}
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' && line[len([]rune(line))-1] == '#' {
			if i > 0 {
				testCases[number][len(testCases[number])-1] = strings.Join(lines[start:i], "\n") + "\n"
			}
			start = i + 1
			number++
		}
	}
}