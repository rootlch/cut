package cut

import (
	"fmt"
	"io"
	"bytes"
)

type Cut []byte

func New() *Cut {
	b := make([]byte, 0)
	h := new(Cut)
	*h = b
	return h
}

func (h *Cut) Write(b []byte) (int, error) {
	*h = append(*h, b...)
	return len(b), nil
}

func (h *Cut) String() string {
	return string(*h)
}

type Data []byte

func (d Data) String() string {
	return string(d)
}

//Paste is a chanel so you can retrieve data from it. But you can also print out
//the data straight away with fmt.Print(Paste)
type Paste <-chan Data

//Data ouput includes the left and right string.
func (h *Cut) Between(left, right string) Paste {
	c := make(chan Data)

	go func() {
		s := *h
		startPat, endPat := []byte(left), []byte(right)
		newline := []byte("\n")
		space := []byte(" ")

		for {
			startPos := bytes.Index(s, startPat)

			if startPos == -1 || len(s) == 0 {
				close(c)
				break
			}

			s = s[startPos:]
			endPos := bytes.Index(s, endPat) + len(endPat)

			c <- bytes.Replace(s[:endPos], newline, space, -1)
			s = s[endPos:]
		}
	}()

	return Paste(c)
}

//This will use up all the tokens in the channel.
func (p Paste) String() string {
	var r string
	for data := range p {
		r += fmt.Sprintln(data)
	}
	return r
}

//Reads the data from the Paste channel to b.
//A "\n" will be used to indicate each data.
//This will use up all the tokens in the channel.
func (p Paste) Read(b []byte) (int, error) {
	data := []byte(<-p)
	if len(data) == 0 {
		return 0, io.EOF
	}
	data = append(data, []byte("\n")...)
	copy(b, data)
	return len(data), nil
}
