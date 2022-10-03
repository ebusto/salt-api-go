package event

import (
	"regexp"

	"github.com/ebusto/salt-api-go"
)

// Event represents a parsed event.
type Event interface{}

type JobNew struct {
	Arguments  []string      `json:"arg"`
	Function   string        `json:"fun"`
	Job        string        `json:"jid"`
	Minions    []string      `json:"minions"`
	Target     salt.Response `json:"tgt"`
	TargetType string        `json:"tgt_type"`
	Time       Time          `json:"_stamp"`
	User       string        `json:"user"`
}

type JobReturn struct {
	Arguments  []string      `json:"fun_args"`
	Command    string        `json:"cmd"`
	Function   string        `json:"fun"`
	Job        string        `json:"jid"`
	Minion     string        `json:"id"`
	Output     string        `json:"out"`
	Return     salt.Response `json:"return"`
	ReturnCode int           `json:"retcode"`
	Success    bool          `json:"success"`
	Time       Time          `json:"_stamp"`
}

type MinionAuth struct {
	Key    string `json:"pub"`
	Minion string `json:"id"`
	Result bool   `json:"result"`
	Status string `json:"act"`
	Time   Time   `json:"_stamp"`
}

type MinionBeacon struct {
	Data   salt.Response `json:"data"`
	Minion string        `json:"id"`
	Name   string        `name:"name"`
	Time   Time          `json:"_stamp"`
}

type MinionKey struct {
	Minion string `json:"id"`
	Result bool   `json:"result"`
	Status string `json:"act"`
	Time   Time   `json:"_stamp"`
}

type MinionRefresh struct {
	Minion string `name:"id"`
	Time   Time   `json:"_stamp"`
}

type MinionStart struct {
	Minion string `json:"id"`
	Time   Time   `json:"_stamp"`
}

type PresenceChange struct {
	Lost []string `json:"lost"`
	New  []string `json:"new"`
	Time Time     `json:"_stamp"`
}

type PresencePresent struct {
	Minions []string `json:"present"`
	Time    Time     `json:"_stamp"`
}

// https://docs.saltstack.com/en/latest/topics/event/master_events.html
var Types = map[*regexp.Regexp]func() Event{
	regexp.MustCompile(`minion/refresh/(?P<id>[^/]+)`):      New[MinionRefresh],
	regexp.MustCompile(`salt/auth`):                         New[MinionAuth],
	regexp.MustCompile(`salt/beacon/[^/]+/(?P<name>[^/]+)`): New[MinionBeacon],
	regexp.MustCompile(`salt/job/\d+/new`):                  New[JobNew],
	regexp.MustCompile(`salt/job/\d+/ret`):                  New[JobReturn],
	regexp.MustCompile(`salt/key`):                          New[MinionKey],
	regexp.MustCompile(`salt/minion/[^/]+/start`):           New[MinionStart],
	regexp.MustCompile(`salt/presence/change`):              New[PresenceChange],
	regexp.MustCompile(`salt/presence/present`):             New[PresencePresent],
}

// New returns a new event of the specified type.
func New[T any]() Event {
	return new(T)
}
