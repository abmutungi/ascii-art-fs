package main

import (
	"fmt"
	"log"
	"os"

	"git.learn.01founders.co/abmutungi/ascii-art-fs.git/am"
)

func main() {
	arg := os.Args

	// error handling
	if len(arg) != 3 {
		fmt.Print("Usage: go run . [STRING] [BANNER]\n")
		fmt.Println()
		fmt.Println("EX: go run . something standard")
		os.Exit(0)
	}

	// defining argument [1] "input string" and [2] "font/banner"
	args := os.Args[1]
	args2 := os.Args[2]
	// This tells it to print a new line if the arg is solely a new line.
	if args == "\\n" {
		fmt.Println()
	} else if args != "" {
		/*The else clause above tells the program to do nothing if
		the argument is an empty string with the rest of the program
		only running if the arg is not an empty string*/

		/* The func splitlines splits the string of the arg into
		a slice of slices split whenever there is a new line*/

		splitLines := am.SplitLines(args)

		lines, err := am.ReadLines(args2 + ".txt")
		if err != nil {
			log.Fatalf("ReadLines: %s", err)
		}

		/*The line below uses the make method to make a map
		and uses a start point of 32 to match up the ascii values
		of each character to the ascii version of the character*/
		charMap := make(map[int][]string)

		start := 32

		for i := 0; i < len(lines); i++ {
			// Tells it to add to start every 9 to match the chars
			if len(charMap[start]) == 9 {
				start++
			}
			charMap[start] = append(charMap[start], lines[i])
		}
		/*The j below refers to the index of each slice within a
		created by splitlines, represented by val. The k represents
		the length of each individual slice. The i iterates up to 9
		to match the height of each character.*/
		for j, val := range splitLines {
			for i := 1; i < 9; i++ {
				for k := 0; k < len(val); k++ {
					fmt.Print(charMap[int(splitLines[j][k])][i])
				}
				fmt.Println()
			}
		}
	}
}
