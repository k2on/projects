package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/k2on/sigma"
	"github.com/rivo/tview"
)

type Conversation struct {
	id    int
	label string
}

type Message struct {
	body string
	from string
}

const CONTACTS_COMMAND = "/usr/local/bin/format-abook"

func getContacts() string {
	cmd := exec.Command(CONTACTS_COMMAND)
	
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil { panic(err) }

	return out.String()
}

func getConversations(client sigma.Client) []Conversation {
	chats, err := client.Chats()
	if err != nil { panic(err) }

	contacts := getContacts()

	getContactName := func (id string) string {
		for _, line := range strings.Split(contacts, "\n") {
			if len(line) == 0 { continue }
			lineParts := strings.Split(line, "=")
			ide := lineParts[0]
			val := lineParts[1]
			if ide == id { return val }
		}
		return id
	}

	conversations := []Conversation{} 
	
	for _, chat := range chats {

		label := getContactName(chat.DisplayName)
		if chat.IsGroupChat {
			label = "(GC) " + label
		}
		conversations = append(conversations, Conversation{chat.ID, label})
	}

	return conversations
}

func getMessagesFromConversation(client sigma.Client, id int) []Message {
	chats, err := client.Messages(id, sigma.MessageFilter{Limit: 50})
	if err != nil { panic(err) }

	handles, err := client.Handles()
	if err != nil { panic(err) }

	getIDFromHandle := func (id string) string {
		for _, handle := range handles {
			if strconv.Itoa(handle.ID) == id { return handle.Identifier }
		}
		return "unknown"
	}

	messages := []Message{}
	for _, chat := range chats {
		var label string
		if chat.FromMe {
			label = "me"
		} else {
			label = getIDFromHandle(chat.Account)
		}
		messages = append([]Message{{chat.Text, label}}, messages... )
	}
	return messages
}

const DEBUG = true

const RUNE_LEFT   = 'h'
const RUNE_DOWN   = 'j'
const RUNE_UP     = 'k'
const RUNE_RIGHT  = 'l'
const RUNE_TOP    = 'g'
const RUNE_BOTTOM = 'G'

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
			
			case event.Rune() == RUNE_TOP:
				fn = func (_ int) int { return 0 }
			
			case event.Rune() == RUNE_BOTTOM:
				fn = func (_ int) int { return endPosition }

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

	// nP := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	// }
	client, err := sigma.NewClient()
	if err != nil { panic(err) }

	app := tview.NewApplication();



	list := tview.NewList()
	conversations := getConversations(client)

	for _, converstation := range conversations {
		list.AddItem(converstation.label, "", 0, nil)
	}


	msg := tview.NewList().
		SetSelectedFocusOnly(true)
		
	
		
	inpt := tview.NewInputField()

	chatId := -1

	inpt.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				msg.SetCurrentItem(-1)
				app.SetFocus(msg)
			} else if key == tcell.KeyEnter {
				message := inpt.GetText()
				msg.AddItem("<me> " + message, "", 0, nil)
				inpt.SetText("")
				msg.SetCurrentItem(-1)

				if !DEBUG {
					client.SendMessage(chatId, message)
				}
			}
		})

	

	addBindings(
		list,
		func(newPos int) {
			chatId = conversations[newPos].id
			msgs := getMessagesFromConversation(client, chatId)
			msg.Clear()
			for _, message := range msgs {
				txt := "<" + message.from + "> " + message.body
				msg.AddItem(txt, "", 0, nil)
			}
			msg.SetCurrentItem(-1)
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
