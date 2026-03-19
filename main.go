package main

import (
	"fmt"
	"time"
)

type Pet struct {
	Name   string
	hunger int
	mood   int
}

func NewPet(name string) *Pet {
	return &Pet{
		Name:   name,
		hunger: 50,
		mood:   50,
	}
}

func (p *Pet) Feed() {
	p.hunger -= 20
	if p.hunger < 0 {
		p.hunger = 0
	}
	p.mood += 5
	if p.mood > 100 {
		p.mood = 100
	}
	fmt.Printf(" %s поел! Голод: %d, Настроение: %d\n", p.Name, p.hunger, p.mood)
}

func (p *Pet) Play() {
	if p.hunger > 70 {
		fmt.Println(" Питомец слишком голоден для игр!")
		return
	}
	p.mood += 20
	if p.mood > 100 {
		p.mood = 100
	}
	p.hunger += 10
	if p.hunger > 100 {
		p.hunger = 100
	}
	fmt.Printf(" %s поиграл! Настроение: %d, Голод: %d\n", p.Name, p.mood, p.hunger)
}

func (p *Pet) Status() {
	emoji := "😐"
	if p.hunger > 80 {
		emoji = "😫"
	} else if p.mood < 30 {
		emoji = "😢"
	} else if p.hunger < 30 && p.mood > 80 {
		emoji = "😍"
	}
	fmt.Printf("🐶 %s %s | Голод: %d | Настроение: %d\n", p.Name, emoji, p.hunger, p.mood)
}

func (p *Pet) Online(stopChan chan bool, alertChan chan string) {
	for {
		select {
		case <-stopChan:
			return
		default:
			time.Sleep(4 * time.Second)
			p.hunger += 5
			p.mood -= 5
			if p.hunger > 100 {
				p.hunger = 100
			}
			if p.mood < 0 {
				p.mood = 0
			}
			if p.hunger >= 90 {
				select {
				case alertChan <- fmt.Sprintf("спасай %s умирает от голода!", p.Name):
				default:
				}
			}
		}
	}
}

func main() {
	var names string
	fmt.Println("привет назови своего питомца")
	fmt.Scan(&names)

	p := NewPet(names)

	stopLife := make(chan bool)
	alerts := make(chan string, 5)

	go p.Online(stopLife, alerts)

	fmt.Printf("%s живет своей жизнью...", names)

	for {
		select {
		case msg := <-alerts:
			fmt.Println("!!!!", msg)
		default:
		}

		fmt.Println("1. кормить | 2. играть | 3. статус | 4. выйти")
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			p.Feed()
		case 2:
			p.Play()
		case 3:
			p.Status()
		case 4:
			stopLife <- true
			time.Sleep(100 * time.Millisecond)
			fmt.Println("пока!")
			return
		default:
			fmt.Println("Ошибка ввода")
		}
	}
}
