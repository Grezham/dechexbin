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

func CreateQuestion(maxRange int, qType int, aType int) question {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	value := r1.Intn(maxRange-1) + 1
	str := ""

	switch qType {
	case binary:
		str = fmt.Sprintf("Convert Binary 0b%08b to ", value)
	case hexadecimal:
		str = fmt.Sprintf("Convert Hexadecimal 0x%x to ", value)
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

func (q *question) Want() string {
	switch q.answerType {
	case binary:
		return fmt.Sprintf("%08b", q.value)
	case hexadecimal:
		return fmt.Sprintf("%x", q.value)
	case decimal:
		return fmt.Sprintf("%d", q.value)
	default:
		return "No Type Found"
	}

}

type SetSettings struct {
	SetSize      int
	MaxRange     int
	QuestionType int
	AnswerType   int
}

type QuestionSet struct {
	questions   []question
	index       int
	results     []bool
	done        bool
	setSettings SetSettings
}

func CreateQuestionSet(setSize int, maxRange int, qType int, aType int) *QuestionSet {
	if setSize >= 1 {
		qSet := make([]question, setSize)
		marks := make([]bool, setSize)

		for i := 0; i < setSize; i++ {
			qSet[i] = CreateQuestion(maxRange, qType, aType)
			marks[i] = false

		}
		return &QuestionSet{
			questions: qSet,
			index:     0,
			results:   marks,
			done:      false,
			setSettings: SetSettings{
				SetSize:      setSize,
				MaxRange:     maxRange,
				QuestionType: qType,
				AnswerType:   aType,
			},
		}

	} else {
		return &QuestionSet{}
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

//I may not need this. Could just format in View() <- main.go
//func (qs *QuestionSet) PrintResults(){}

func (qs QuestionSet) isDone() bool {
	return qs.done
}

func (qs *QuestionSet) Reset() {
	qs.index = 0
	qs.done = false
	qs.results = make([]bool, len(qs.results))
}

func (qs *QuestionSet) Restart() {
	qs.Reset()
	sSize := qs.setSettings.SetSize
	tmpq := make([]question, sSize)
	for i := 0; i < sSize; i++ {
		tmpq[i] = CreateQuestion(qs.setSettings.MaxRange, qs.setSettings.QuestionType, qs.setSettings.AnswerType)

	}
	qs.questions = tmpq
}
