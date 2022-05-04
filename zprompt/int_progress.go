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

func (i *intProgress) Set(val int) error {
	i.Value = val
	return i.Draw(false)
}

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

func (i *intProgress) Complete() error {
	i.Value = i.Total
	return i.Draw(true)
}
