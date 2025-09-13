package input

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/eiannone/keyboard"
	"github.com/manifoldco/promptui"
)

type ValidationFunc func(input string) error

type SearcherFunc func(input string, index int) bool

type IInputMethods interface {
	StringInput(result *string, label string, validation *ValidationFunc) error
	MultipleStringInput(result *[]string, label string, validation *ValidationFunc) error
	IntInput(result *int, label string, validation *ValidationFunc) error
	Int64Input(result *int64, label string, validation *ValidationFunc) error
	UInt32Input(result *uint32, label string, validation *ValidationFunc) error
	UInt64Input(result *uint64, label string, validation *ValidationFunc) error
	Float32Input(result *float32, label string, validation *ValidationFunc) error
	Float64Input(result *float64, label string, validation *ValidationFunc) error
	BooleanInput(result *bool, label string) error
	SelectInput(result *string, label string, items []string) error
	MultipleSelectInput(result *[]string, label string, items []string, searchFunc *SearcherFunc) error
}

type InputMethods struct{}

func NewInputMethods() *InputMethods {
	return &InputMethods{}
}

func (im *InputMethods) StringInput(result *string, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	*result = input
	return nil
}

func (im *InputMethods) MultipleStringInput(result *[]string, label string, validation *ValidationFunc) error {
	if result == nil {
		return errors.New("result ptr is nil")
	}

	var entries []string
	for {
		var entry string
		if err := im.StringInput(&entry, label+" (or leave blank to finish)", validation); err != nil {
			return err
		}
		if entry == "" {
			break
		}
		entries = append(entries, entry)
	}
	*result = entries
	return nil
}

func (im *InputMethods) IntInput(result *int, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value int
	value, err = strconv.Atoi(input)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) UInt32Input(result *uint32, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value uint64
	value, err = strconv.ParseUint(input, 10, 32)
	if err != nil {
		return err
	}

	*result = uint32(value)
	return nil
}

func (im *InputMethods) Int64Input(result *int64, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value int64
	value, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) UInt64Input(result *uint64, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value uint64
	value, err = strconv.ParseUint(input, 10, 64)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) Float32Input(result *float32, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value float64
	value, err = strconv.ParseFloat(input, 32)
	if err != nil {
		return err
	}

	*result = float32(value)
	return nil
}

func (im *InputMethods) Float64Input(result *float64, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation == nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value float64
	value, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) BooleanInput(result *bool, label string) error {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"true", "false"},
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "{{ . | cyan }}",
			Inactive: "{{ . | faint }}",
			Selected: "{{ . | green }}",
		},
	}

	_, input, err := prompt.Run()
	if err != nil {
		return err
	}

	value, err := strconv.ParseBool(input)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) SelectInput(result *string, label string, items []string) error {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "{{ . | cyan }}",
			Inactive: "{{ . | faint }}",
			Selected: "{{ . | green }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return err
	}

	if index < 0 || index >= len(items) {
		return errors.New("invalid selection index")
	}

	*result = items[index]
	return nil
}

func (im *InputMethods) MultipleSelectInput(result *[]string, label string, items []string, searchFunc *SearcherFunc) error {
	if result == nil {
		return errors.New("result ptr is nil")
	}
	if len(items) == 0 {
		*result = []string{}
		return nil
	}

	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	activeStyle := promptui.Styler(promptui.FGCyan)
	selectedStyle := promptui.Styler(promptui.FGGreen)
	labelStyle := promptui.Styler(promptui.FGBold)
	faintStyle := promptui.Styler(promptui.FGFaint)

	selected := map[int]struct{}{}
	cursor := 0
	query := ""

	type entry struct {
		idx   int
		label string
	}

	filter := func() []entry {
		if query == "" && searchFunc == nil {
			out := make([]entry, 0, len(items))
			for i, v := range items {
				out = append(out, entry{idx: i, label: v})
			}
			return out
		}
		lowerQ := strings.ToLower(query)
		out := []entry{}
		for i, v := range items {
			ok := true
			if searchFunc != nil {
				ok = (*searchFunc)(query, i)
			} else if query != "" {
				ok = strings.Contains(strings.ToLower(v), lowerQ)
			}
			if ok {
				out = append(out, entry{idx: i, label: v})
			}
		}
		return out
	}

	clearScreen := func() { fmt.Print("\033[H\033[2J") }

	render := func(listing []entry) {
		clearScreen()
		fmt.Println(labelStyle(fmt.Sprintf("%s (↑/↓: Move Space: Select Enter: Confirm / ESC: Cancel)", label)))
		if query != "" {
			fmt.Println(faintStyle("検索:") + " " + query)
		}
		for i, e := range listing {
			cur := (i == cursor)
			_, isSel := selected[e.idx]
			check := "[ ]"
			if isSel {
				check = "[x]"
			}
			display := e.label
			if cur && isSel {
				display = selectedStyle(display)
			} else if cur {
				display = activeStyle(display)
			} else if isSel {
				display = selectedStyle(display)
			}
			pointer := "  "
			if cur {
				pointer = "> "
			}
			fmt.Printf("%s%s %s\n", pointer, check, display)
		}
		if len(listing) == 0 {
			fmt.Println(faintStyle("Not Found") + " (Backspace to clear search)")
		}
	}

	listing := filter()
	render(listing)

	for {
		r, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}
		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return errors.New("canceled")
		case keyboard.KeyEnter:
			// finalize
			out := make([]string, 0, len(selected))
			for i, v := range items {
				if _, ok := selected[i]; ok {
					out = append(out, v)
				}
			}
			*result = out
			return nil
		case keyboard.KeyArrowUp:
			if len(listing) > 0 {
				cursor--
				if cursor < 0 {
					cursor = len(listing) - 1
				}
			}
		case keyboard.KeyArrowDown:
			if len(listing) > 0 {
				cursor++
				if cursor >= len(listing) {
					cursor = 0
				}
			}
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			if query != "" {
				query = query[:len(query)-1]
				listing = filter()
				if cursor >= len(listing) {
					cursor = len(listing) - 1
				}
			}
		case keyboard.KeySpace: // toggle selection
			if len(listing) > 0 {
				idx := listing[cursor].idx
				if _, ok := selected[idx]; ok {
					delete(selected, idx)
				} else {
					selected[idx] = struct{}{}
				}
			}
		default:
			if unicode.IsPrint(r) {
				query += string(r)
				listing = filter()
				if cursor >= len(listing) {
					cursor = len(listing) - 1
				}
			}
		}
		render(listing)
	}
}
