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

type Client interface {
	GetContactName(identifier string) string

	GetConversations() ([]Conversation, error)

	GetConversationMessages(id int) ([]Message, error)
}

type realClient struct {
	config Config
	contacts string
	sigmaClient sigma.Client
	handles []sigma.Handle
}

func PipeCommand(in string, command string) string {
	cmd := exec.Command("/bin/sh", "-c", "echo " + in + " | " + command)
	
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil { panic(err) }

	return out.String()
}

func RunCommand(command string) string {
	cmd := exec.Command(command)
	
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil { panic(err) }

	return out.String()
}

type Config struct {
	contacts string
	nameMessage string
}

func NewClient(config Config) (Client, error) {
	contacts := RunCommand(config.contacts)
	sigmaClient, err := sigma.NewClient()
	if err != nil { panic(err) }
	handles, err := sigmaClient.Handles()
	if err != nil { panic(err) }

	return &realClient{
		config,
		contacts,
		sigmaClient,
		handles,
	}, nil
}

func (client *realClient) GetContactName(phoneOrEmail string) string {
	contacts := client.contacts
	lines := strings.Split(contacts, "\n")
	for _, line := range lines {
		isLineEmpty := len(line) == 0
		if isLineEmpty { continue }
		lineParts := strings.Split(line, "=")
		id := lineParts[0]
		value := lineParts[1]
		if phoneOrEmail == id { return value }
	}
	return phoneOrEmail
}

func (client *realClient) GetIdFromHandle(handle string) string {
	handles := client.handles
	for _, handleItem := range handles {
		if strconv.Itoa(handleItem.ID) == handle { return handleItem.Identifier }
	}
	return "unknown"
}

func (client *realClient) GetConversations() ([]Conversation, error) {
	chats, err := client.sigmaClient.Chats()
	if err != nil { return []Conversation{}, err }

	conversations := []Conversation{} 

	for _, chat := range chats {
		label := client.GetContactName(chat.DisplayName)
		if chat.IsGroupChat { label = "(GC) " + label }
		conversations = append(conversations, Conversation{chat.ID, label})
	}
	return conversations, nil
}

func (client *realClient) GetConversationMessages(id int) ([]Message, error) {
	messagesRaw, err := client.sigmaClient.Messages(id, sigma.MessageFilter{Limit: 50})
	if err != nil { return []Message{}, err }

	var cachedNames = make(map[string]string)

	messages := []Message{}
	for _, chat := range messagesRaw {
		var label string
		if chat.FromMe {
			label = "我"
		} else {
			id := client.GetIdFromHandle(chat.Account)
			label = client.GetContactName(id) 
			name, exists := cachedNames[label]
			if exists {
				label = name
			} else {
				name = PipeCommand(label, client.config.nameMessage)
				cachedNames[label] = name
				label = name
			}
		}
		messages = append([]Message{{chat.Text, label}}, messages... )
	}
	return messages, nil
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

	config := Config{
		contacts: "/usr/local/bin/format-abook",
		nameMessage: "awk '{print $1;}'",
	}
	client, err := NewClient(config)
	if err != nil { panic(err) }

	app := tview.NewApplication();



	list := tview.NewList()
	conversations, err := client.GetConversations()

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
					// client.SendMessage(chatId, message)
				}
			}
		})

	

	addBindings(
		list,
		func(newPos int) {
			chatId = conversations[newPos].id
			msgs, err := client.GetConversationMessages(chatId)
			if err != nil { panic(err) }

			msg.Clear()

			longestName := 0

			for _, msg := range msgs {
				length := len(msg.from)
				if length > longestName { longestName = length }
			}

			for _, message := range msgs {
				nameSizeDiff :=  longestName - (len(message.from))
				padding := strings.Repeat(" ", nameSizeDiff)

				txt := padding + message.from + " : " + message.body
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
