package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: strings.TrimSpace(line[0]),
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func shuffleProblems(problemsList *[]problem) {
	problems := *problemsList
	for i := range problems {
		j := rand.Intn(i + 1)
		problems[i], problems[j] = problems[j], problems[i]
	}
}

func fetchQuizDataFromCSV(csvFilename *string) ([][]string, error) {

	//_ = csvFilename
	// do ./quizGame --help to find all the variables the binery is using

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	return r.ReadAll()

}
func scanAnswer(answerCh chan<- string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	answer = strings.TrimSpace(answer)
	answerCh <- answer
}
func doQuiz(problems []problem, timeLimit *int) {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	answerCh := make(chan string)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		go scanAnswer(answerCh)
		select {
		case <-timer.C:
			fmt.Printf("\nYou got %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
				fmt.Println("Correct")
			}
		}
	}
	close(answerCh)
	fmt.Printf("You got %d out of %d\n", correct, len(problems))
}
func setFlags() (*string, *int) {
	csvFilename := flag.String("csv", "problems.csv",
		"csv in the format for 'question,answer'")
	timeLimit := flag.Int("time", 30, "time to take the quize")
	flag.Parse()
	return csvFilename, timeLimit
}
func main() {
	csvFilename, timeLimit := setFlags()
	lines, err := fetchQuizDataFromCSV(csvFilename)
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)
	fmt.Println(problems)

	shuffleProblems(&problems)

	doQuiz(problems, timeLimit)

}
