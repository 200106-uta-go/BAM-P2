// Package userinputs is a collections of functions for requesting
// information from the user
package userinputs

import (
	"fmt"
	"strconv"
	"strings"
)

// RequestAnswer takes a string (question), displays it to the user and returns
// the response in the form of a string (answer)
func RequestAnswer(question string) string {
	var a string
	fmt.Println(question) // output question
	fmt.Scan(&a)          // get answer
	return a              // return answer
}

// MultiChoiceAnswer takes a string (question) and an array of string (answers)
// displays both to user and returns the response in the form of a string (answer)
// it compaires the (answer) to (answers), it will re-ask the (question) and
// redisplay the (answers) if a valid (answer) is not found. (answer) can be a match
// to the answer or the numberic refrence displayed (index of answers + 1) but not a
// combination of both.
func MultiChoiceAnswer(question string, answers []string) string {
	var answer string
	answered := false    // is questioned answered correctly
	answeredNum := false // is answer numeric
	for !answered {
		// display question
		fmt.Println(question)
		// display answers array
		for i := 0; i < len(answers); i++ {
			q := i + 1                 // index + 1 (rl represintation of order)
			fmt.Println(q, answers[i]) // display each answer posible
		}

		fmt.Scan(&answer)

		// check answer
		for i := 0; i < len(answers); i++ {
			if strings.ToLower(answers[i]) == strings.ToLower(answer) {
				answered = true
				break
			} else if answer == strconv.Itoa(i+1) {
				answered = true
				answeredNum = true
				break
			}
		}

		// failed
		if !answered {
			fmt.Println("Invalid answer, try again...")
		}
	}

	if answeredNum {
		answerNum, _ := strconv.Atoi(answer)
		answerNum--
		answer = answers[answerNum]
		return answer
	}

	return answer
}
