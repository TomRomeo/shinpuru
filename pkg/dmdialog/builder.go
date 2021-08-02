package dmdialog

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Result int

const (
	ResultOK Result = iota
	ResultCanceled
	ResultTimeout
)

type InputValidator func(string) error

type Question struct {
	ID        string
	Text      string
	Validator InputValidator
}

type Builder struct {
	session         *discordgo.Session
	timeout         time.Duration
	questions       []*Question
	chanID          string
	cursor          int
	currentQuestion *Question
}

func New(s *discordgo.Session) *Builder {
	return &Builder{
		session:   s,
		questions: make([]*Question, 0),
		timeout:   5 * time.Minute,
	}
}

func (b *Builder) WithTimeout(d time.Duration) *Builder {
	b.timeout = d
	return b
}

func (b *Builder) AddQuestion(q Question) *Builder {
	if q.ID == "" {
		q.ID = strconv.Itoa(len(b.questions))
	}
	if q.Validator == nil {
		q.Validator = func(s string) error {
			return nil
		}
	}
	b.questions = append(b.questions, &q)
	return b
}

func (b *Builder) Send(userID string) (a *Awaiter, err error) {
	if len(b.questions) == 0 {
		err = errors.New("no questions were set")
		return
	}

	ch, err := b.session.UserChannelCreate(userID)
	if err != nil {
		return
	}

	b.chanID = ch.ID

	a = &Awaiter{
		answers:   make(map[string]string),
		cFinished: make(chan Result, 1),
	}

	go func() {
		time.Sleep(b.timeout)
		a.setResult(ResultTimeout)
		b.session.ChannelMessageSend(b.chanID, "⌛️ Timed out.")
	}()

	if _, err = b.sendNextQuestion(); err != nil {
		return
	}

	a.removeHandler = b.session.AddHandler(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		if e.ChannelID != ch.ID || e.Author.ID != userID {
			return
		}

		switch strings.ToLower(e.Content) {
		case "cancel", "exit", "stop", "quit":
			a.setResult(ResultCanceled)
			b.session.ChannelMessageSend(b.chanID, "🛑 Canceled.")
			return
		}

		if _, ok := a.answers[b.currentQuestion.ID]; ok {
			return
		}

		if err := b.currentQuestion.Validator(e.Content); err != nil {
			b.session.ChannelMessageSend(b.chanID, "⚠️ Input Error: "+err.Error()+"\n\nPlease enter your anwer again.")
			return
		}

		a.answers[b.currentQuestion.ID] = e.Content

		if ok, err := b.sendNextQuestion(); !ok && err == nil {
			a.setResult(ResultOK)
		}
	})

	return
}

func (b *Builder) sendNextQuestion() (ok bool, err error) {
	if b.cursor == len(b.questions) {
		return
	}

	q := b.questions[b.cursor]
	if _, err = b.session.ChannelMessageSend(b.chanID, q.Text); err != nil {
		return
	}

	b.currentQuestion = q
	b.cursor++
	ok = true

	return
}
