// Package prompt provides user interaction utilities for CLI applications.
package prompt

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const defaultCancelMessage = "⚠️  Operation canceled"

var (
	// ErrInvalidSelection indicates the user entered an out-of-range selection.
	ErrInvalidSelection = errors.New("invalid selection")
	// ErrInvalidConfirmation indicates the response could not be parsed as yes/no.
	ErrInvalidConfirmation = errors.New("invalid confirmation")
)

// Interface defines the prompt operations used by the CLI commands.
type Interface interface {
	Input(prompt string) (string, bool, error)
	Select(title string, items []string, prompt string) (int, bool, error)
	Confirm(prompt string) (bool, bool, error)
}

// Prompter reads from an input stream and writes prompts/results to an output stream.
type Prompter struct {
	baseReader    io.Reader
	reader        *bufio.Reader
	writer        io.Writer
	cancelMessage string
	inputFile     *os.File
}

type inputOutcome struct {
	line     string
	canceled bool
	err      error
}

// NewDefault returns a prompter wired to stdin/stdout.
func NewDefault() Interface {
	return New(os.Stdin, os.Stdout)
}

// New creates a prompter backed by the provided reader and writer.
func New(reader io.Reader, writer io.Writer) Interface {
	if reader == nil {
		reader = os.Stdin
	}
	if writer == nil {
		writer = os.Stdout
	}

	p := &Prompter{
		baseReader:    reader,
		reader:        bufio.NewReader(reader),
		writer:        writer,
		cancelMessage: defaultCancelMessage,
	}
	if file, ok := reader.(*os.File); ok {
		p.inputFile = file
	}
	return p
}

// WithCancelMessage allows overriding the cancel message used by the prompter.
func (p *Prompter) WithCancelMessage(message string) *Prompter {
	if message != "" {
		p.cancelMessage = message
	}
	return p
}

// Input prompts the user for free-form input and returns the response.
// The returned bool is true when the input was canceled via soft cancel.
func (p *Prompter) Input(prompt string) (string, bool, error) {
	if p == nil {
		return "", true, nil
	}

	if outcome, used := p.inputWithTerminal(prompt); used {
		return outcome.line, outcome.canceled, outcome.err
	}

	return p.inputWithBuffered(prompt)
}

func (p *Prompter) inputWithBuffered(prompt string) (string, bool, error) {
	if _, err := fmt.Fprint(p.writer, prompt); err != nil {
		return "", false, err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	line, err := p.reader.ReadString('\n')
	if ctx.Err() != nil {
		p.printCancelMessage()
		return "", true, nil
	}
	if err != nil {
		if errors.Is(err, io.EOF) {
			if line == "" {
				return "", false, io.EOF
			}
			return trimNewline(line), false, nil
		}
		if isInterrupted(err) {
			p.printCancelMessage()
			return "", true, nil
		}
		return "", false, err
	}

	return trimNewline(line), false, nil
}

func (p *Prompter) inputWithTerminal(prompt string) (inputOutcome, bool) {
	terminal, restore, ok := p.prepareTerminal()
	if !ok {
		return inputOutcome{}, false
	}
	defer restore()

	terminal.SetPrompt(prompt)

	line, readErr := terminal.ReadLine()
	switch {
	case readErr == nil:
		return inputOutcome{line: line}, true
	case errors.Is(readErr, term.ErrPasteIndicator):
		return inputOutcome{line: line}, true
	case errors.Is(readErr, io.EOF):
		if line != "" {
			return inputOutcome{line: line}, true
		}
		p.printCancelMessage()
		return inputOutcome{canceled: true}, true
	default:
		return inputOutcome{err: readErr}, true
	}
}

func (p *Prompter) prepareTerminal() (*term.Terminal, func(), bool) {
	if p == nil || p.inputFile == nil {
		return nil, nil, false
	}

	fd := int(p.inputFile.Fd())
	if !term.IsTerminal(fd) {
		return nil, nil, false
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, nil, false
	}

	restore := func() {
		_ = term.Restore(fd, state)
		if p.reader != nil && p.baseReader != nil {
			p.reader.Reset(p.baseReader)
		}
	}

	terminal := term.NewTerminal(&terminalReadWriter{
		reader: p.inputFile,
		writer: p.writer,
	}, "")
	return terminal, restore, true
}

// Select displays a numbered list and prompts the user to choose an item.
// It returns a zero-based index on success.
func (p *Prompter) Select(title string, items []string, prompt string) (int, bool, error) {
	if p == nil {
		return -1, true, nil
	}
	if len(items) == 0 {
		return -1, false, fmt.Errorf("no items to select")
	}

	if err := p.displaySelection(title, items); err != nil {
		return -1, false, err
	}

	line, canceled, err := p.Input(prompt)
	if canceled || err != nil {
		return -1, canceled, err
	}

	return p.parseSelection(line, len(items))
}

func (p *Prompter) displaySelection(title string, items []string) error {
	if title != "" {
		if _, err := fmt.Fprintln(p.writer, title); err != nil {
			return err
		}
	}
	for i, item := range items {
		if _, err := fmt.Fprintf(p.writer, "[%d] %s\n", i+1, item); err != nil {
			return err
		}
	}
	return nil
}

func (p *Prompter) parseSelection(line string, itemCount int) (int, bool, error) {
	idx, convErr := strconv.Atoi(strings.TrimSpace(line))
	if convErr != nil || idx < 1 || idx > itemCount {
		return -1, false, ErrInvalidSelection
	}
	return idx - 1, false, nil
}

// Confirm prompts the user for a yes/no answer; defaults to "no" on empty input.
func (p *Prompter) Confirm(prompt string) (bool, bool, error) {
	if p == nil {
		return false, true, nil
	}

	line, canceled, err := p.Input(prompt)
	if canceled || err != nil {
		return false, canceled, err
	}

	normalized := strings.ToLower(strings.TrimSpace(line))
	switch normalized {
	case "", "n", "no":
		return false, false, nil
	case "y", "yes":
		return true, false, nil
	default:
		return false, false, fmt.Errorf("%w: %s", ErrInvalidConfirmation, line)
	}
}

func (p *Prompter) printCancelMessage() {
	if p == nil || p.writer == nil {
		return
	}
	if _, err := fmt.Fprintf(p.writer, "\n%s\n", p.cancelMessage); err != nil {
		_ = err // explicitly ignore printing errors
	}
}

type terminalReadWriter struct {
	reader *os.File
	writer io.Writer
}

func (rw *terminalReadWriter) Read(p []byte) (int, error) {
	if rw == nil || rw.reader == nil {
		return 0, io.EOF
	}
	return rw.reader.Read(p)
}

func (rw *terminalReadWriter) Write(p []byte) (int, error) {
	if rw == nil || rw.writer == nil {
		return 0, io.EOF
	}
	return rw.writer.Write(p)
}

func trimNewline(s string) string {
	s = strings.TrimSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\r")
	return s
}

func isInterrupted(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, syscall.EINTR) {
		return true
	}
	// Some platforms wrap EINTR differently; fall back to substring match.
	return strings.Contains(strings.ToLower(err.Error()), "interrupt")
}
