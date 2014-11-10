package weezard_test

import (
  "log"
  "github.com/aliafshar/weezard"
)

func Example() {

  type UserInfo struct {
    Name string `question:",What is your name?"`
    DateOfBirth string `question:",What is your date of birth?"`
  }

  u := &UserInfo{}
  err := weezard.Ask(u)
  if err != nil {
    log.Fatalln(err)
  }
  log.Println(u)
}

func main() {
  Example()
}