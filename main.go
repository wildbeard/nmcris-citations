package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func Contains(a []int, b int) bool {
	for _, i := range a {
		if i == b {
			return true
		}
	}
	return false
}

func main() {
	f, err := os.Open("citations.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	output := [][]string{}
	// Title, Author, Year, Report #, NMCRIS #, # Sites, Acreage, IALLRESOUR, ISURVEYACR
	// Title (Author Year) (Report No. #) (NMCRIS No. #)
	toTake := []int{0, 2, 3, 6, 9, 12, 13, 14, 15}

	for i, line := range lines {
		if i == 0 {
			continue
		}

		lineData := []string{}

		for j, field := range line {
			if !Contains(toTake, j) {
				continue
			}
			switch j {
			case 0:
				lineData = append(lineData, strings.Trim(field, " "))
			case 2:
				lName := strings.Split(field, ",")[0]
				if strings.Index(field, "And") != -1 {
					lName = lName + " et al"
				}
				lineData = append(lineData, fmt.Sprintf("(%v", lName))
			case 3:
				lineData = append(lineData, fmt.Sprintf("%v)", field))
			case 6:
				lineData = append(lineData, fmt.Sprintf("(Report No. %v)", field))
			case 9:
				lineData = append(lineData, fmt.Sprintf("(NMCRIS No. %v)", field))
			default:
				lineData = append(lineData, field)
			}
		}

		output = append(output, lineData)
		//fmt.Printf("%v\n", lineData)
	}

	outFile, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, os.ModePerm)

	if err != nil {
		log.Fatalf("Could not open output file: %v\n", err)
	}

	defer outFile.Close()

	for _, line := range output {
		outLine := ""
		for j, field := range line {
			if j != 5 {
				outLine = outLine + " " + field
			} else if j == 5 {
				outLine = outLine + "\n" + field
			}
		}

		_, err := outFile.WriteString(outLine + "\n\n")

		if err != nil {
			log.Fatalf("Could not write line: %v\n", err)
		}
	}

}
