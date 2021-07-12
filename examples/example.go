package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ebusto/salt-api-go"
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
	c.Minions.All(ctx, printMinion)

	// Display a list of minions matching '*salt*'.
	c.Minions.Filter(ctx, "*salt*", printMinion)

	// Display a list of all jobs.
	c.Jobs.All(ctx, printJob)

	// Ping all minions.
	c.Ping(ctx, "*", printPong)

	// Collect disk usage on Ubuntu 18.04 minions with a 10 second timeout.
	cmd := salt.Command{
		Client:     "local",
		Function:   "disk.usage",
		Target:     "osfinger:Ubuntu-18.04",
		TargetType: "grain",
		Timeout:    10,
	}

	c.Run(ctx, &cmd, printUsage)

	// Fire an event onto the bus.
	c.Events.Fire(ctx, "salt-api-go/test", salt.Request{
		"test": true,
		"time": time.Now(),
	})

	// Display events for 10 seconds.
	ctx, _ = context.WithTimeout(ctx, time.Second*10)

	c.Events.Stream(ctx, printEvent)
}

func printMinion(id string, grains salt.Response) error {
	log.Println(id, grains.Get("osfinger"), grains.Get("saltversion"))

	return nil
}

func printJob(id string, job salt.Response) error {
	log.Println(id, job.Get("User"), job.Get("Target"), job.Get("Function"), job.Get("Arguments"))

	return nil
}

func printPong(id string, ok bool) error {
	log.Println(id, ok)

	return nil
}

func printUsage(id string, response salt.Response) error {
	// The return from Salt for disk.uage unfortunately encodes all values as
	// strings, even if they're integers, hence the `json:"<field>,string"`.
	var usage map[string]struct {
		Available  int    `json:"available,string"`
		Capacity   string `json:"capacity"`
		Filesystem string `json:"filesystem"`
		Used       int    `json:"used,string"`
	}

	if err := response.Decode(&usage); err != nil {
		log.Fatal(err)
	}

	for mount, info := range usage {
		log.Println(id, mount, info.Filesystem, info.Used)
	}

	return nil
}

func printEvent(event salt.Response) error {
	log.Println(event.Get("tag"), string(event))

	return nil
}
