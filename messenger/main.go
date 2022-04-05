package main

import (
	// "fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// client, err := sigma.NewClient()
// if err != nil { panic(err) }

// chats, err := client.Chats()
// if err != nil { panic(err) }

// chats, err := client.Messages(22, sigma.MessageFilter{Limit: 50})
// if err != nil { panic(err) }

// for _, chat := range chats {
// fmt.Println(chat)
// }

const RUNE_LEFT  = 'h'
const RUNE_DOWN  = 'j'
const RUNE_UP    = 'k'
const RUNE_RIGHT = 'l'

func addBindings(list *tview.List, hoverFn func(i int), leftFn func(i int), rightFn func(i int)) {
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			endPosition := list.GetItemCount() - 1
			position := list.GetCurrentItem()
			var fn func(i int) int
			switch {
			case event.Rune() == RUNE_DOWN:
				fn = func(i int) int {
					if i == endPosition { return 0 }
					return i + 1
				}
			case event.Rune() == RUNE_UP:
				fn = func(i int) int {
					if i == 0 { return endPosition }
					return i - 1
				}
			case event.Rune() == RUNE_RIGHT:
				rightFn(position)
				return nil
			
			case event.Rune() == RUNE_LEFT:
				leftFn(position)
				return nil

			default:
				return nil
			}
			
			newPos := fn(position)
			list.SetCurrentItem(newPos)
			hoverFn(newPos)
			
			return event
		})

}


func main() {

	m := make(map[int][]string)
	m[0] = []string{"mes 1", "mes 2"}
	m[1] = []string{"mes 3", "mes 4"}

	// nP := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	// }

	app := tview.NewApplication();

	list := tview.NewList().
		AddItem("List item 1", "", 0, nil).
		AddItem("List item 2", "", 0, nil).
		AddItem("List item 3", "", 0, nil).
		AddItem("List item 4", "", 0, nil)

	msg := tview.NewList().
		AddItem("my message", "", 0, nil).
		SetSelectedFocusOnly(true)
		
	inpt := tview.NewInputField()

	addBindings(
		list,
		func(newPos int) {
			msgs := m[newPos]
			msg.Clear()
			for _, message := range msgs {
				msg.AddItem(message, "", 0, nil)
			}
		},
		func(_ int) {
			app.Stop()
		},
		func(i int) {
			app.SetFocus(inpt)			
		},
	)
	 
	addBindings(
		msg,
		func(_ int) {},
		func(_ int) {
			app.SetFocus(list)
		},
		func(_ int) {},
	)
	 
		
	msgr := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		SetBorders(true).
		AddItem(msg, 0, 0, 1, 1, 0, 0, false).
		AddItem(inpt, 1, 0, 1, 1, 0, 0, false)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(list, 0, 0, 1, 1, 0, 0, true).
		AddItem(msgr, 0, 1, 1, 1, 0, 0, false)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}


}
