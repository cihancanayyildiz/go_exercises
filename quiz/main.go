package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type quiz struct {
	question string
	answer string
}

func newQuiz(question string, answer string) *quiz {
	return &quiz{question, answer}
}

var quizSlc []quiz

func main() {
	filename := flag.String("file", "problems.csv", "Questions csv file")
	timelimit := flag.Int("time", 30, "Time limit variable")

	flag.Parse()

	fd, error := os.Open(*filename)

	if error != nil {
		fmt.Println(error)
	}

	reader := csv.NewReader(fd)

	records, error := reader.ReadAll()

	if error != nil {
		fmt.Println(error)
	}

	for _, record := range records {
		quizSlc = append(quizSlc, *newQuiz(record[0], record[1]))
	}

	right := 0

	// 30 seconds
	timer := time.NewTimer(time.Second * time.Duration(*timelimit))

	go func() {
		<-timer.C
		fmt.Println("Time is up!")
		fmt.Printf("total number of questions correct %d\n",right)
		fmt.Printf("total number of questions %d\n",len(quizSlc))
		os.Exit(0)
	}()

	var answer string
	for _, quizItem := range quizSlc {
		fmt.Println(quizItem.question)
		fmt.Scanf("%s", &answer)
		if answer == quizItem.answer {
			right++
		}
	}

	defer fd.Close()
}