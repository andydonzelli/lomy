package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const REGION_ID = "my_messages_region"

type Tui struct {
	InputFieldQueue <-chan string
	app             *tview.Application
	textView        *tview.TextView
	inputField      *tview.InputField
}

func CreateTui() Tui {
	app := tview.NewApplication()

	textView := createTextViewChatHistory(app)
	inputField := createInputField()

	inputFieldQueue := make(chan string, 3)
	configureInputFieldFunc(inputField, inputFieldQueue, textView)

	flex := createFlex(textView, inputField)

	app.SetRoot(flex, true).SetFocus(inputField)

	return Tui{
		app:             app,
		textView:        textView,
		inputField:      inputField,
		InputFieldQueue: inputFieldQueue,
	}
}

func (tui *Tui) RunAppAndBlock() {
	if err := tui.app.Run(); err != nil {
		panic(err)
	}
}

func (tui *Tui) WriteToTextView(str string) {
	fmt.Fprint(tui.textView, str)
}

func createTextViewChatHistory(app *tview.Application) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		Highlight(REGION_ID).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		}).SetText(" ... Welcome to the chat ...\n\n")
	textView.
		SetBorder(true).SetTitle(" Lomy Chat ")

	return textView
}

func createInputField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).SetFieldBackgroundColor(tcell.ColorBlack)
}

func configureInputFieldFunc(inputField *tview.InputField, inputFieldQueue chan string, textView *tview.TextView) {
	const REGION_START = "[\"" + REGION_ID + "\"]"
	const REGION_END = "[\"\"]"
	inputField.
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				inputFieldQueue <- inputField.GetText() + "\n"
				fmt.Fprintf(textView, "%sYou: %s%s\n",
					REGION_START, inputField.GetText(), REGION_END)
				inputField.SetText("")
			}
		})
}

func createFlex(textView *tview.TextView, inputField *tview.InputField) *tview.Flex {
	return tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(textView, 0, 1, false).
			AddItem(inputField, 2, 1, true), 0, 2,
			false)
}
