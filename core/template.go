package core

import (
	"bytes"
	"text/template"

	color "github.com/aybabtme/rgbterm"
	"github.com/mgutz/ansi"
)

var DisableColor = false

var (
	HelpInputRune = '?'

	ErrorIcon    = "✘"
	HelpIcon     = "ⓘ"
	QuestionIcon = "?"

	MarkedOptionIcon   = "◉"
	UnmarkedOptionIcon = "◯"

	SelectFocusIcon = "❯"
)

func red(s string) string {
	return color.FgString(s, 194, 37, 92)
}

func purple(s string) string {
	return color.FgString(s, 96, 97, 190)
}

var TemplateFuncs = map[string]interface{}{
	// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
	"color": func(color string) string {
		if DisableColor {
			return ""
		}

		switch color {
		case "default":
			return "\033[38;5;61m"
		case "red":
			return "\033[38;5;125m"
		}

		return ansi.ColorCode(color)
	},
	"HelpInputRune": func() string {
		return string(HelpInputRune)
	},
	"ErrorIcon": func() string {
		return ErrorIcon
	},
	"HelpIcon": func() string {
		return HelpIcon
	},
	"QuestionIcon": func() string {
		return QuestionIcon
	},
	"MarkedOptionIcon": func() string {
		return MarkedOptionIcon
	},
	"UnmarkedOptionIcon": func() string {
		return UnmarkedOptionIcon
	},
	"SelectFocusIcon": func() string {
		return SelectFocusIcon
	},
}

var memoizedGetTemplate = map[string]*template.Template{}

func getTemplate(tmpl string) (*template.Template, error) {
	if t, ok := memoizedGetTemplate[tmpl]; ok {
		return t, nil
	}

	t, err := template.New("prompt").Funcs(TemplateFuncs).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	memoizedGetTemplate[tmpl] = t
	return t, nil
}

func RunTemplate(tmpl string, data interface{}) (string, error) {
	t, err := getTemplate(tmpl)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
