package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ebusto/salt-api-go"
	"github.com/ebusto/salt-api-go/event"
)

func main() {
	c := salt.New(os.Getenv("SALTAPI_URL"))

	ctx := context.Background()

	err := c.Login(ctx,
		os.Getenv("SALTAPI_USER"),
		os.Getenv("SALTAPI_PASS"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer c.Logout(ctx)

	// Display a list of all minions.
	c.Minions.All(ctx, func(id string, grains salt.Response) error {
		log.Printf("Minion %s, os %s", id, grains.Get("osfinger"))

		return nil
	})

	// Display a list of all jobs.
	c.Jobs.All(ctx, func(id string, job salt.Response) error {
		log.Printf("Job %s, user %s, target %s, function %s", id,
			job.Get("User"),
			job.Get("Target"),
			job.Get("Function"),
		)

		return nil
	})

	// Display a list of all accepted minion keys.
	c.Keys.ListAccepted(ctx, func(id string) error {
		log.Printf("Minion %s", id)

		return nil
	})

	// Ping all minions.
	c.Ping(ctx, "*", func(id string, ok bool) error {
		log.Printf("Minion %s, pong %t", id, ok)

		return nil
	})

	// Collect disk usage on Ubuntu 18.04 minions with a 10 second timeout.
	var cmd = salt.Command{
		Client:     "local",
		Function:   "disk.usage",
		Target:     "osfinger:Ubuntu-18.04",
		TargetType: "grain",
		Timeout:    10,
	}

	// The return for disk.uage unfortunately encodes all values as strings.
	var usage map[string]struct {
		Available  int    `json:"available,string"`
		Capacity   string `json:"capacity"`
		Filesystem string `json:"filesystem"`
		Used       int    `json:"used,string"`
	}

	c.Run(ctx, &cmd, func(id string, response salt.Response) error {
		if err := response.Decode(&usage); err != nil {
			log.Fatal(err)
		}

		for mount, info := range usage {
			log.Println(id, mount, info.Filesystem, info.Used)
		}

		return nil
	})

	// Fire an event onto the bus.
	c.Events.Fire(ctx, "salt-api-go/test", salt.Request{
		"test": true,
		"time": time.Now(),
	})

	// Display events for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	p := event.NewParser()

	c.Events.Stream(ctx, func(response salt.Response) error {
		event, err := p.Parse(response)

		if err != nil {
			return err
		}

		// A nil event indicates an unhandled event.
		if event != nil {
			log.Printf("%#v", event)
		}

		return nil
	})
}
