//
//Package cut provide utilities for tokenizing data.
package cut

/*
Copyright (c) 2012 Lim Chia Hau. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither my name. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"bytes"
	"fmt"
	"io"
)

type Cut []byte

//This is the first function your need to call before doing anything else.
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
