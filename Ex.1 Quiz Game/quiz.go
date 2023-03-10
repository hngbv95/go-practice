package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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

func getAnswer() (<-chan string, <-chan error) {
	// Since we are comunicating using channel
	// Error would be returned using channel too
	answChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		var answer string
		_, err := fmt.Scanf("%s\n", &answer)

		// Only return err/answer, not both
		if err != nil {
			errorChan <- err
		} else {
			answChan <- strings.ToLower(strings.TrimSpace(answer))
		}
		//Close right after use
		close(answChan)
		close(errorChan)
	}()

	return answChan, errorChan
}

func startQuiz(quizList []Quiz, timeout <-chan time.Time) (int, error) {
	// Print Questions and Read user input, collect the answer
	correctAnswer := 0
	for _, quiz := range quizList {
		fmt.Printf("Question : %v = ", quiz.question)
		answerChan, errorChan := getAnswer()

		select {
		case <-timeout:
			fmt.Println("\nStop! Time up!")
			return correctAnswer, nil
		case answer := <-answerChan:
			if answer == quiz.answer {
				correctAnswer++
			}
		case err := <-errorChan:
			return correctAnswer, err
		}
	}

	return correctAnswer, nil
}

func ReadFlag() (string, int) {
	// Read arguments flag from command line
	filePath := flag.String("f", "problem.csv", "Specify a file path for questions (.csv file)")
	duration := flag.Int("t", 30, "The test duration")
	// suffle := flag.Bool("s", false, "should be suffle or not")
	flag.Parse()

	return *filePath, *duration
}

func main() {

	filePath, duration := ReadFlag()

	quizList, err := getQuiz(filePath)
	if err != nil {
		log.Fatal(err)
	}

	//Set Timer
	timer := time.NewTimer(time.Duration(duration) * time.Second)

	correctAnswer, err := startQuiz(quizList, timer.C)
	if err != nil {
		log.Fatal(err)
	}

	// Return the statistic
	fmt.Printf("\nYou are correct %v/%v questions", correctAnswer, len(quizList))
}
