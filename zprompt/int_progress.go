package zprompt

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func NewProgress(fd *os.File, total int) *intProgress {
	i := int(fd.Fd())
	return &intProgress{
		Total:         total,
		ProgressChar:  "=",
		DrawFrequency: time.Second,
		terminalId:    i,
		fd:            fd,
	}
}

type intProgress struct {
	Total         int
	Value         int
	ProgressChar  string
	DrawFrequency time.Duration
	Message       string

	fd            *os.File
	lastDraw      time.Time
	lastDrawValue int
	terminalId    int
}

// Set the value of progress bar to a value
func (i *intProgress) Set(val int) error {
	i.Value = val
	return i.Draw(false)
}

// Set a message with max length
func (i *intProgress) SetMessage(msg string, length int) {
	center := "..."
	if len(msg) > length {
		toRemove := len(msg) + len(center) - length
		index := len(msg) / 2
		i.Message = msg[:index-toRemove/2] + center + msg[index+toRemove/2:]
	}
}

// Increment the progress bar by one
func (i *intProgress) Inc() error {
	return i.Set(i.Value + 1)
}

// Draw the progress bar on the terminal. This will exit early if it is drawing too early
func (i *intProgress) Draw(force bool) (err error) {
	w, _, err := terminal.GetSize(i.terminalId)
	if err != nil {
		return
	}
	now := time.Now()
	if !force && now.Before(i.lastDraw.Add(i.DrawFrequency)) {
		return
	}
	i.lastDraw = now
	i.lastDrawValue = i.Value
	p := (float64(i.Value) / float64(i.Total))
	header := fmt.Sprintf("%.0f%% ", p*100)
	trailer := fmt.Sprintf(" %d/%d %s", i.Value, i.Total, i.Message)
	w = w - len(trailer) - len(header) - 1
	nBar := int(float64(w) * p)
	bar := strings.Repeat(i.ProgressChar, nBar)
	nSpace := w - nBar
	empty := ""
	if nSpace > 0 {
		empty = strings.Repeat(" ", nSpace)
	}
	_, err = fmt.Fprintf(i.fd, "\r%s%s%s%s", header, bar, empty, trailer)
	return
}

// Set the progress bar to 100% complete and force a draw operation
func (i *intProgress) Complete() (err error) {
	i.Value = i.Total
	i.Message = ""
	if err := i.Draw(true); err != nil {
		return err
	}
	_, err = fmt.Fprintln(i.fd)
	return
}

func (i *intProgress) Clear() (err error) {
	_, err = fmt.Fprintf(i.fd, "\r")
	return
}

// Reset the progress bar
func (i *intProgress) Reset(total int) {
	i.Value = 0
	i.Total = total
}
