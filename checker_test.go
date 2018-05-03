package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecker(t *testing.T) {
	_, err := newChecker(rootURL, 1, false)
	assert.Nil(t, err)
}

func TestNewCheckerError(t *testing.T) {
	_, err := newChecker(":", 1, false)
	assert.NotNil(t, err)
}

func TestCheckerCheck(t *testing.T) {
	for _, s := range []string{rootURL, fragmentURL, baseURL, redirectURL} {
		c, _ := newChecker(s, 1, false)

		go c.Check()

		for r := range c.Results() {
			assert.True(t, r.OK())
		}
	}
}

func TestCheckerCheckPage(t *testing.T) {
	c, _ := newChecker(rootURL, 256, false)

	r, err := c.fetcher.FetchLink(existentURL)
	assert.Nil(t, err)

	p, ok := r.Page()
	assert.True(t, ok)

	go c.checkPage(p)

	assert.True(t, (<-c.Results()).OK())
}

func TestCheckerCheckPageError(t *testing.T) {
	for _, s := range []string{erroneousURL, invalidBaseURL} {
		c, _ := newChecker(rootURL, 256, false)

		r, err := c.fetcher.FetchLink(s)
		assert.Nil(t, err)

		p, ok := r.Page()
		assert.True(t, ok)

		go c.checkPage(p)

		assert.False(t, (<-c.Results()).OK())
	}
}

func TestStringChannelToSlice(t *testing.T) {
	for _, c := range []struct {
		channel chan string
		slice   []string
	}{
		{
			make(chan string, 1),
			[]string{},
		},
		{
			func() chan string {
				c := make(chan string, 1)
				c <- "foo"
				return c
			}(),
			[]string{"foo"},
		},
		{
			func() chan string {
				c := make(chan string, 2)
				c <- "foo"
				c <- "bar"
				return c
			}(),
			[]string{"foo", "bar"},
		},
		{
			func() chan string {
				c := make(chan string, 3)
				c <- "foo"
				c <- "bar"
				c <- "baz"
				return c
			}(),
			[]string{"foo", "bar", "baz"},
		},
	} {
		assert.Equal(t, c.slice, stringChannelToSlice(c.channel))
	}
}
