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
	1:  {"hello", ""},
	2:  {"HELLO", ""},
	3:  {"HeLlo HuMaN", ""},
	4:  {"1Hello 2There", ""},
	5:  {"Hello\\nThere", ""},
	6:  {"{Hello & There #}", ""},
	7:  {"hello There 1 to 2!", ""},
	8:  {"MaD3IrA&LiSboN", ""},
	9:  {"1a\"#FdwHywR&/()=", ""},
	10: {"{|}~", ""},
	11: {"[\\]^_ 'a", ""},
	12: {"RGB", ""},
	13: {":;<=>?@", ""},
	14: {"\\!\" #$%&'()*+,-./", ""},
	15: {"ABCDEFGHIJKLMNOPQRSTUVWXYZ", ""},
	16: {"abcdefghijklmnopqrstuvwxyz", ""},
}

/*	This test file tests the ascii-art project against the first 16 test cases on
	audit page	*/
func TestAsciiArt(t *testing.T) {
	getTestCases()

	/*	Iterate through each test case and starting a goroutine for each, this
		is done so instead of waiting for the previous test to complete they can
		all be checked simulaneously	*/
	var wg sync.WaitGroup
	for i := 1; i <= len(testCases); i++ {
		wg.Add(1)
		go func(current []string, w *sync.WaitGroup, ti *testing.T) {
			defer w.Done()
			result := getResult(current)
			/*	Fails the project if the test cases expected output doesn't match
				the actual output	*/
			if string(result) != current[1] {
				ti.Errorf("\nTest fails when given the test case:\n\t\"%s\","+
					"\nexpected:\n%s\ngot:\n%s\n\n",
					current[0], current[1], string(result))
			}
		}(testCases[i], &wg, t)
	}
	wg.Wait()
}

/*	This function imitates the running of "go run . string", which it then pipes
	into a second function "cat -e" to immitate and then returns the result	*/
func getResult(testCase []string) string {
	first := exec.Command("go", "run", ".", testCase[0])
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
			testCases[number][1] = strings.Join(lines[start:], "\n") + "\n"
			break
		}
		if line[0] == '#' && line[len([]rune(line))-1] == '#' {
			if i > 0 {
				testCases[number][1] = strings.Join(lines[start:i], "\n") + "\n"
			}
			start = i + 1
			number++
		}
	}
}