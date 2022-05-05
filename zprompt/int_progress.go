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
	}
}

type intProgress struct {
	Total         int
	Value         int
	ProgressChar  string
	DrawFrequency time.Duration

	lastDraw   time.Time
	terminalId int
}

// Set the value of progress bar to a value
func (i *intProgress) Set(val int) error {
	i.Value = val
	return i.Draw(false)
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
	if !force && time.Now().Before(i.lastDraw.Add(i.DrawFrequency)) {
		return
	}
	trailer := fmt.Sprintf(" %d/%d", i.Value, i.Total)
	w = w - len(trailer) - 1
	nBar := int(float64(w) * (float64(i.Value) / float64(i.Total)))
	bar := strings.Repeat(i.ProgressChar, nBar)
	nSpace := w - nBar
	space := ""
	if nSpace > 0 {
		space = strings.Repeat(" ", nSpace)
	}
	fmt.Printf("\r%s%s%s", bar, space, trailer)
	return
}

// Set the progress bar to 100% complete and force a draw operation
func (i *intProgress) Complete() error {
	i.Value = i.Total
	return i.Draw(true)
}

// Reset the progress bar
func (i *intProgress) Reset() error {
	i.Value = 0
	return i.Draw(true)
}
