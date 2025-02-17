package survey

import (
	"errors"

	"github.com/h2non/survey/core"
)

// PageSize is the default maximum number of items to show in select/multiselect prompts
var PageSize = 7

// Validator is a function passed to a Question after a user has provided a response.
// If the function returns an error, then the user will be prompted again for another
// response.
type Validator func(interface{}) error

// Question is the core data structure for a survey questionnaire.
type Question struct {
	Name     string
	Prompt   Prompt
	Validate Validator
}

// Prompt is the primary interface for the objects that can take user input
// and return a response.
type Prompt interface {
	Prompt() (interface{}, error)
	Cleanup(interface{}) error
	Error(error) error
}

/*
AskOne performs the prompt for a single prompt and asks for validation if required.
Response types should be something that can be casted from the response type designated
in the documentation. For example:

	name := ""
	prompt := &survey.Input{
		Message: "name",
	}

	survey.AskOne(prompt, &name, nil)

*/
func AskOne(p Prompt, response interface{}, v Validator) error {
	err := Ask([]*Question{{Prompt: p, Validate: v}}, response)
	if err != nil {
		return err
	}

	return nil
}

/*
Ask performs the prompt loop, asking for validation when appropriate. The response
type can be one of two options. If a struct is passed, the answer will be written to
the field whose name matches the Name field on the corresponding question. Field types
should be something that can be casted from the response type designated in the
documentation. Note, a survey tag can also be used to identify a Otherwise, a
map[string]interface{} can be passed, responses will be written to the key with the
matching name. For example:

	qs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "What is your name?"},
			Validate: survey.Required,
		},
	}

	answers := struct{ Name string }{}


	err := survey.Ask(qs, &answers)
*/
func Ask(qs []*Question, response interface{}) error {

	// if we weren't passed a place to record the answers
	if response == nil {
		// we can't go any further
		return errors.New("cannot call Ask() with a nil reference to record the answers")
	}

	// go over every question
	for _, q := range qs {
		// grab the user input and save it
		ans, err := q.Prompt.Prompt()
		// if there was a problem
		if err != nil {
			return err
		}

		// if there is a validate handler for this question
		if q.Validate != nil {
			// wait for a valid response
			for invalid := q.Validate(ans); invalid != nil; invalid = q.Validate(ans) {
				err := q.Prompt.Error(invalid)
				// if there was a problem
				if err != nil {
					return err
				}

				// ask for more input
				ans, err = q.Prompt.Prompt()
				// if there was a problem
				if err != nil {
					return err
				}
			}
		}

		// tell the prompt to cleanup with the validated value
		q.Prompt.Cleanup(ans)

		// if something went wrong
		if err != nil {
			// stop listening
			return err
		}

		// add it to the map
		err = core.WriteAnswer(response, q.Name, ans)
		// if something went wrong
		if err != nil {
			return err
		}

	}
	// return the response
	return nil
}

// paginate returns a single page of choices given the page size, the total list of
// possible choices, and the current selected index in the total list.
func paginate(page int, choices []string, sel int) ([]string, int) {
	// the number of elements to show in a single page
	var pageSize int
	// if the select has a specific page size
	if page != 0 {
		// use the specified one
		pageSize = page
		// otherwise the select does not have a page size
	} else {
		// use the package default
		pageSize = PageSize
	}

	var start, end, cursor int

	if len(choices) < pageSize {
		// if we dont have enough options to fill a page
		start = 0
		end = len(choices)
		cursor = sel

	} else if sel < pageSize/2 {
		// if we are in the first half page
		start = 0
		end = pageSize
		cursor = sel

	} else if len(choices)-sel-1 < pageSize/2 {
		// if we are in the last half page
		start = len(choices) - pageSize
		end = len(choices)
		cursor = sel - start

	} else {
		// somewhere in the middle
		above := pageSize / 2
		below := pageSize - above

		cursor = pageSize / 2
		start = sel - above
		end = sel + below
	}

	// return the subset we care about and the index
	return choices[start:end], cursor
}
