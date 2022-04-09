package main

import (
	"bytes"
	"fmt"
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
	id          string
	parentId    string
	messageType string
	body        string
	from        string
	reaction    string
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
	contacts         string
	nameMessage      string
	reactionLove     string
	reactionLike     string
	reactionDislike  string
	reactionLaugh    string
	reactionEmphasis string
	reactionQuestion string
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
		if chat.Type != 0 { 
			reaction := "?"
			switch {
				case chat.Type == REACTION_ADD_LOVE:
					reaction = client.config.reactionLove
				case chat.Type == REACTION_ADD_LIKE:
					reaction = client.config.reactionLike
				case chat.Type == REACTION_ADD_DISLIKE:
					reaction = client.config.reactionDislike
				case chat.Type == REACTION_ADD_LAUGH:
					reaction = client.config.reactionLaugh
				case chat.Type == REACTION_ADD_EMPHASIS:
					reaction = client.config.reactionEmphasis
				case chat.Type == REACTION_ADD_QUESTION:
					reaction = client.config.reactionQuestion
			}
			messages = append([]Message{{
				chat.Guid,
				chat.ParentId,
				"reaction",
				"",
				"",
				reaction,
			}}, messages...)
			continue
		}
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
		messages = append([]Message{{
			chat.Guid,
			"",
			"text",
			chat.Text,
			label,
			"",
		}}, messages... )
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

const REACTION_ADD_LOVE     = 2000
const REACTION_ADD_LIKE     = 2001
const REACTION_ADD_DISLIKE  = 2002
const REACTION_ADD_LAUGH    = 2003
const REACTION_ADD_EMPHASIS = 2004
const REACTION_ADD_QUESTION = 2005

const REACTION_REMOVE_LOVE     = 3000
const REACTION_REMOVE_LIKE     = 3001
const REACTION_REMOVE_DISLIKE  = 3002
const REACTION_REMOVE_LAUGH    = 3003
const REACTION_REMOVE_EMPHASIS = 3004
const REACTION_REMOVE_QUESTION = 3005

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
func Map(vs []string, f func(string) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

func reactionsFormat(reactions []Message) string {
	if len(reactions) == 0 { return "" }

	var reactionsFormatted []string

	for _, reaction := range reactions {
		reactionsFormatted = append(reactionsFormatted, reaction.reaction)
	}

	emojis := strings.Join(reactionsFormatted, " ")

	return fmt.Sprintf("[%s]", emojis)
	
}

func main() {

	config := Config{
		contacts: "/usr/local/bin/format-abook",
		nameMessage: "awk '{print $1;}'",
		reactionLove: "💕",
		reactionLike: "👍",
		reactionDislike: "👎",
		reactionLaugh: "🤣",
		reactionEmphasis: "❗",
		reactionQuestion: "❓",
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
			var childrenMap = make(map[string][]Message)

			for _, msg := range msgs {
				length := len(msg.from)
				if length > longestName { longestName = length }

				// parentId := msg.parentId
				parentId := strings.Replace(msg.parentId, "p:0/", "", 1)
				// fmt.Println(parentId)
				if parentId != "" {
					_, exists := childrenMap[parentId]
					if !exists {
						childrenMap[parentId] = []Message{}
					}
					childrenMap[parentId] = append(childrenMap[parentId], msg)

				}


			}



			for _, message := range msgs {
				if message.messageType != "text" { continue }

				nameSizeDiff :=  longestName - (len(message.from))
				padding := strings.Repeat(" ", nameSizeDiff)

				// fmt.Println(message.id)
				reactions := reactionsFormat(childrenMap[message.id])

				txt := padding + message.from + " : " + message.body + " " + reactions
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
