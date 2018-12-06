package main

import (
	"fmt"
	"flag"
	"github.com/260by/tools/translate"
)

func main()  {
	en := flag.String("en", "", "english to chinese")
	flag.Parse()

	ch, err := translate.GetTranslateEnToChContent(*en)
	if err != nil {
		panic(err)
	}

	fmt.Println(ch)
}