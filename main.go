package main

import (
	"fmt"

	"github.com/anthonyminyungi/learngo/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	baseWord := "hello"
	dictionary.Add(baseWord, "First")
	found, _ := dictionary.Search(baseWord)
	fmt.Println(found)
	err := dictionary.Delete(baseWord)
	if err != nil {
		fmt.Println(err)
	}
	word, err2 := dictionary.Search(baseWord)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(word)
}
