# weezard

Populates golang structs with questions on the command line. The library takes
its metadata from struct field tags in the same way you would for `json`, etc.

**Reference docs:** http://godoc.org/github.com/aliafshar/weezard

An example will make it more clear, perhaps.

    package weezard_test

    import (
      "log"
      "github.com/aliafshar/weezard"
    )

    func Example() {
      type UserInfo struct {
        Name string `question:"default,What is your name?"`
        DateOfBirth string `question:",What is your date of birth?"`
      }

      u := &UserInfo{}
      err := weezard.Ask(u)
      if err != nil {
        log.Fatalln(err)
      }
      log.Println(u)
    }

The tag format is:

    question:"<default>,<question>"
