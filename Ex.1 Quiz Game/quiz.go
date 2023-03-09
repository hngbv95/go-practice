package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Quiz struct {
	question string
	answer   string
}

func getQuiz(fileName string) ([]Quiz, error) {
	// Read file and parse csv to struct object

	// Be cafeful Readfile will read entire file, so large file may corrupt your system
	// Should we change to file stream instead
	file, err := os.ReadFile(fileName)
	// return error if can't open
	if err != nil {
		return nil, err
	}

	// Parse csv to object
	// we have a reader
	r := csv.NewReader(strings.NewReader(string(file)))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	// check the length before access element
	//fmt.Println(records[0][0])
	quizList := []Quiz{}
	for i := 0; i < len(records); i++ {
		quiz := Quiz{records[i][0], records[i][1]}
		quizList = append(quizList, quiz)
	}

	return quizList, nil
}

func startQuiz(quizList []Quiz) (int, error) {
	// Print Questions and Read user input, collect the answer
	correctAnswer := 0
	for _, quiz := range quizList {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Printf("Question : %v\nYour Answer: \n", quiz.question)
		//read input
		answer, err := inputReader.ReadString('\n')
		if err != nil {
			return correctAnswer, err
		}

		answer = strings.TrimSpace(answer)
		fmt.Println()

		if answer == quiz.answer {
			correctAnswer++
		}
	}

	return correctAnswer, nil
}

func main() {
	// Read arguments flag from command line
	filePath := flag.String("f", "problem.csv", "Specify a file path for questions (.csv file)")
	duration := flag.Int("t", 30, "The test duration")
	suffle := flag.Bool("s", false, "should be suffle or not")
	flag.Parse()
	fmt.Println(*filePath)

	quizList, err := getQuiz(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	correctAnswer, err := startQuiz(quizList)
	if err != nil {
		log.Fatal(err)
	}

	// Return the statistic
	fmt.Printf("You are correct %v/%v questions\n", correctAnswer, len(quizList))
}
