package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Card struct that represents a cards against humanity card
type Card struct {
	CardBody string
	CardType int64
	Blanks   int64
	Classic  bool
}

func getScanner() *bufio.Scanner {
	return bufio.NewScanner(os.Stdin)
}

func getContinueInput(scanner *bufio.Scanner) bool {
	ok := scanner.Scan()
	if !ok {
		return false
	}
	return scanner.Text() == "y"
}

func createCard(scanner *bufio.Scanner) Card {
	var card Card
	var err error
	fmt.Println("Enter 0 for white card, 1 for black card. Press enter for 0.")
	ok := scanner.Scan()
	if !ok {
		panic("There was a problem scanning for input.")
	}

	inp := scanner.Text()
	if inp == "" {
		card.CardType = 0
	} else {
		card.CardType, err = strconv.ParseInt(inp, 0, 0)
		if err != nil {
			card.CardType = 0
		}
	}

	fmt.Println("Enter card body.")
	ok = scanner.Scan()
	if !ok {
		panic("There was a problem scanning for input.")
	}

	card.CardBody = scanner.Text()
	if card.CardType == 1 {
		fmt.Println("How many blank spaces does the card have?")
		ok = scanner.Scan()
		if !ok {
			panic("There was a problem scanning for input.")
		}

		card.Blanks, err = strconv.ParseInt(scanner.Text(), 0, 0)
		if err != nil {
			panic("There was a problem converting that to a number.")
		}
	} else {
		card.Blanks = 0
	}

	fmt.Println(card)
	return card
}

func writeOutput(cards *[]Card, id int) {
	fileFlags := os.O_APPEND | os.O_RDWR | os.O_CREATE
	fmt.Println("Writing Output.")
	f, err := os.OpenFile("cards_data.txt", fileFlags, os.ModeAppend)
	defer f.Close()
	if err != nil {
		fmt.Println("Could not open the file. Printing cards to screen.")
		fmt.Println(cards)
		panic(err)
	}

	for _, card := range *cards {
		_, err := f.WriteString(fmt.Sprintf("%d,\"%s\",%d,%d\n", id, card.CardBody, card.CardType, card.Blanks))
		if err != nil {
			fmt.Println("Could not write to file. Printing cards to screen.")
			fmt.Println(cards)
			panic(err)
		}
		id++
	}
}

func makeCards() {
	var cards []Card
	cont := true
	inputScanner := getScanner()

	defer func() {
		recover()
		id, err := strconv.ParseInt(os.Args[2], 0, 0)
		if err != nil {
			id = 0
		}
		writeOutput(&cards, int(id))
	}()

	for cont {
		cards = append(cards, createCard(inputScanner))
		fmt.Print("Continue? (y)")
		cont = getContinueInput(inputScanner)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Try running it with:")
		fmt.Println("cardsData input [id to start at]")
		fmt.Println("cardsData write [path to csv file]")
		fmt.Println("cardsData post 'body' [type] [blanks]")
		os.Exit(0)
	}

	if os.Args[1] == "input" {
		makeCards()
	} else if os.Args[1] == "write" {
		WriteToDB()
	} else if os.Args[1] == "post" {
		PostToDB()
	}
}
