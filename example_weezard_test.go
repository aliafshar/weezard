package weezard_test

import (
	"github.com/aliafshar/weezard"
	"log"
)

func ExampleAsk() {
	type UserInfo struct {
		Name        string `question:",What is your name?"`
		DateOfBirth string `question:",What is your date of birth?"`
	}
	u := &UserInfo{}
	err := weezard.Ask(u)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(u)
}

func ExampleAskQuestion() {
  q := &weezard.Question{Usage: "What is your username?"}
  a, err := weezard.AskQuestion(q)
  if err != nil {
    log.Fatalln(err)
  }
  log.Println(a)
}

func ExampleAskQuestion_callback() {
  var a string
  q := &weezard.Question{Usage: "What is your username?", Set: func(v string) { a = v }}
  _, err := weezard.AskQuestion(q)
  if err != nil {
    log.Fatalln(err)
  }
  log.Println(a)
}
