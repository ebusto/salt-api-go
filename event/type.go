package event

import (
	"regexp"

	"github.com/ebusto/salt-api-go"
)

// Event represents a parsed event.
type Event any

// https://docs.saltstack.com/en/latest/topics/event/master_events.html
var Types = map[*regexp.Regexp]func() Event{
	regexp.MustCompile(`minion/refresh/(?P<id>[^/]+)`):      New[MinionRefresh],
	regexp.MustCompile(`salt/auth`):                         New[MinionAuth],
	regexp.MustCompile(`salt/beacon/[^/]+/(?P<name>[^/]+)`): New[MinionBeacon],
	regexp.MustCompile(`salt/job/\d+/new`):                  New[JobNew],
	regexp.MustCompile(`salt/job/\d+/ret`):                  New[JobReturn],
	regexp.MustCompile(`salt/job/\d+/prog/[^/]+/\d+`):       New[JobProgress],
	regexp.MustCompile(`salt/key`):                          New[MinionKey],
	regexp.MustCompile(`salt/minion/[^/]+/start`):           New[MinionStart],
	regexp.MustCompile(`salt/presence/change`):              New[PresenceChange],
	regexp.MustCompile(`salt/presence/present`):             New[PresencePresent],
}

// New returns a new event of the specified type.
func New[T any]() Event {
	return new(T)
}

type JobNew struct {
	Arguments  []any         `json:"arg"`
	Function   string        `json:"fun"`
	Job        string        `json:"jid"`
	Minions    []string      `json:"minions"`
	Target     salt.Response `json:"tgt"`
	TargetType string        `json:"tgt_type"`
	Time       Time          `json:"_stamp"`
	User       string        `json:"user"`
}

type JobReturn struct {
	Arguments  []any         `json:"fun_args"`
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

type JobProgress struct {
	Time   Time   `json:"_stamp"`
	Master string `json:"_master"`
	Cmd    string `json:"cmd"`
	ID     string `json:"id"`
	Job    string `json:"jid"`
	Data   struct {
		Len    int `json:"len"`
		Return struct {
			RunNum    int           `json:"__run_num__"`
			SLS       string        `json:"__sls__"`
			StateID   string        `json:"__id__"`
			Changes   salt.Response `json:"changes"`
			Comment   string        `json:"comment"`
			Duration  float64       `json:"duration"`
			Result    bool          `json:"result"`
			StartTime string        `json:"start_time"`
			Name      string        `json:"name"`
		} `json:"ret"`
	} `json:"data"`
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
