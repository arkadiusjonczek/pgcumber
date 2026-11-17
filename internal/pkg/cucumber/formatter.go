package cucumber

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/formatters"
	messages "github.com/cucumber/messages/go/v34"
)

type stepResult struct {
	text string
	icon string
	err  error
}

type PrettyIconFormatter struct {
	mu           sync.Mutex
	out          io.Writer
	stepKeywords map[string]string

	currentScenario string
	currentSteps    []stepResult

	passed    int
	failed    int
	skipped   int
	undefined int
}

func NewPrettyIconFormatter(_ string, out io.Writer) godog.Formatter {
	return &PrettyIconFormatter{
		out:          out,
		stepKeywords: make(map[string]string),
	}
}

func (f *PrettyIconFormatter) TestRunStarted() {}

func (f *PrettyIconFormatter) Feature(gherkinDocument *messages.GherkinDocument, _ string, _ []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.flushScenario()

	if gherkinDocument.Feature == nil {
		return
	}

	for _, child := range gherkinDocument.Feature.Children {
		if child.Background != nil {
			for _, s := range child.Background.Steps {
				f.stepKeywords[s.Id] = strings.TrimSpace(s.Keyword)
			}
		}
		if child.Scenario != nil {
			for _, s := range child.Scenario.Steps {
				f.stepKeywords[s.Id] = strings.TrimSpace(s.Keyword)
			}
		}
		if child.Rule != nil {
			for _, ruleChild := range child.Rule.Children {
				if ruleChild.Background != nil {
					for _, s := range ruleChild.Background.Steps {
						f.stepKeywords[s.Id] = strings.TrimSpace(s.Keyword)
					}
				}
				if ruleChild.Scenario != nil {
					for _, s := range ruleChild.Scenario.Steps {
						f.stepKeywords[s.Id] = strings.TrimSpace(s.Keyword)
					}
				}
			}
		}
	}

	fmt.Fprintf(f.out, "\nFeature: %s\n", gherkinDocument.Feature.Name)
}

func (f *PrettyIconFormatter) getKeyword(step *messages.PickleStep) string {
	for _, id := range step.AstNodeIds {
		if kw, ok := f.stepKeywords[id]; ok {
			return kw
		}
	}
	return ""
}

// flushScenario gibt das gesammelte Szenario mit bündig ausgerichteten Icons aus
func (f *PrettyIconFormatter) flushScenario() {
	if f.currentScenario == "" && len(f.currentSteps) == 0 {
		return
	}

	fmt.Fprintf(f.out, "\n  Scenario: %s\n", f.currentScenario)

	// Maximale Zeichenlänge (Rune-Count für Umlaute/Sonderzeichen) ermitteln
	maxLen := 0
	for _, s := range f.currentSteps {
		length := utf8.RuneCountInString(s.text)
		if length > maxLen {
			maxLen = length
		}
	}

	// Alle Steps mit Leerzeichen bis maxLen + 2 Spalten Abstand auffüllen
	for _, s := range f.currentSteps {
		padding := strings.Repeat(" ", maxLen-utf8.RuneCountInString(s.text)+2)
		fmt.Fprintf(f.out, "    %s%s%s\n", s.text, padding, s.icon)
		if s.err != nil {
			fmt.Fprintf(f.out, "      Error: %v\n", s.err)
		}
	}

	f.currentScenario = ""
	f.currentSteps = nil
}

func (f *PrettyIconFormatter) Pickle(pickle *messages.Pickle) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.flushScenario()
	f.currentScenario = pickle.Name
}

func (f *PrettyIconFormatter) Defined(_ *messages.Pickle, _ *messages.PickleStep, _ *formatters.StepDefinition) {
}

func (f *PrettyIconFormatter) addStep(step *messages.PickleStep, icon string, err error) {
	line := fmt.Sprintf("%s %s", f.getKeyword(step), step.Text)
	f.currentSteps = append(f.currentSteps, stepResult{
		text: line,
		icon: icon,
		err:  err,
	})
}

func (f *PrettyIconFormatter) Passed(_ *messages.Pickle, step *messages.PickleStep, _ *formatters.StepDefinition) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.passed++
	f.addStep(step, "✅", nil)
}

func (f *PrettyIconFormatter) Failed(_ *messages.Pickle, step *messages.PickleStep, _ *formatters.StepDefinition, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failed++
	f.addStep(step, "❌", err)
}

func (f *PrettyIconFormatter) Skipped(_ *messages.Pickle, step *messages.PickleStep, _ *formatters.StepDefinition) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.skipped++
	f.addStep(step, "⏭️", nil)
}

func (f *PrettyIconFormatter) Undefined(_ *messages.Pickle, step *messages.PickleStep, _ *formatters.StepDefinition) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.undefined++
	f.addStep(step, "❓", nil)
}

func (f *PrettyIconFormatter) Pending(_ *messages.Pickle, step *messages.PickleStep, _ *formatters.StepDefinition) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.undefined++
	f.addStep(step, "🚧", nil)
}

func (f *PrettyIconFormatter) Ambiguous(_ *messages.Pickle, step *messages.PickleStep, _ *formatters.StepDefinition, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failed++
	f.addStep(step, "⚠️", err)
}

func (f *PrettyIconFormatter) Summary() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.flushScenario()

	total := f.passed + f.failed + f.skipped + f.undefined
	fmt.Fprintf(f.out, "\n%d steps (%d passed, %d failed, %d skipped, %d undefined)\n",
		total, f.passed, f.failed, f.skipped, f.undefined)
}
