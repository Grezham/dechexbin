package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type question struct {
	str        string
	value      int
	answer     string
	answerType int
}

func CreateQuestion(qType int, aType int, max int) question {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	value := r1.Intn(max-1) + 1
	str := ""

	switch qType {
	case binary:
		str = fmt.Sprintf("Convert Binary %08b to ", value)
	case hexadecimal:
		str = fmt.Sprintf("Convert Hexadecimal %x to ", value)
	case decimal:
		str = fmt.Sprintf("Convert Decimal %d to ", value)
	}

	switch aType {
	case binary:
		str += "Binary"
	case hexadecimal:
		str += "Hexadecimal"
	case decimal:
		str += "Decimal"
	}

	return question{
		str:        str,
		value:      value,
		answer:     "",
		answerType: aType,
	}
}

type QuestionSet struct {
	questions []question
	index     int
	results   []bool
	done      bool
}

func CreateQuestionSet(setSize int, qType int, aType int, maxRange int) QuestionSet {
	if setSize >= 1 {
		qSet := make([]question, setSize)
		marks := make([]bool, setSize)

		for i := 0; i < setSize; i++ {
			qSet[i] = CreateQuestion(qType, aType, maxRange)
			marks[i] = false

		}
		return QuestionSet{
			questions: qSet,
			results:   marks,
		}

	} else {
		return QuestionSet{}
	}
}

func (qs *QuestionSet) GetCurrentQuestion() string {
	if qs.isDone() {
		return "Done!"
	}

	return qs.questions[qs.index].str
}

func (qs *QuestionSet) NextQuestion() {
	qs.index++
	if qs.index < len(qs.questions) {
		return
	} else {
		qs.done = true
		return
	}
}

func (qs QuestionSet) GetQuestionNumber() int {
	return qs.index + 1
}

func (qs *QuestionSet) GetAnswer(ans string) {
	qs.questions[qs.index].answer = ans
}

func (qs *QuestionSet) CheckAnswer() {
	currQuestion := qs.questions[qs.index]
	num, err := strconv.ParseInt(currQuestion.answer, currQuestion.answerType, 64)
	if err != nil {
		currQuestion.answer += fmt.Sprintf(" %s", err.Error())
		return
	}

	if num == int64(currQuestion.value) {
		qs.results[qs.index] = true
	}
}

func (qs QuestionSet) isDone() bool {
	return qs.done
}

func (qs *QuestionSet) Reset() {
	qs.index = 0
	qs.done = false
	qs.results = make([]bool, len(qs.results))
}
