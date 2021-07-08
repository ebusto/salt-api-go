package main

import (
	"context"
	"log"
	"os"

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

	// Display a list of all minions.
	c.Minions.All(ctx, printMinion)

	// Display a list of minions matching '*salt*'.
	c.Minions.Filter(ctx, "*salt*", printMinion)

	// Display a list of all jobs.
	c.Jobs.All(ctx, printJob)

	c.Logout(ctx)
}

func printMinion(id string, grains salt.Response) error {
	log.Println(id, grains.Get("osfinger"), grains.Get("saltversion"))

	return nil
}

func printJob(id string, job salt.Response) error {
	log.Println(id, job.Get("User"), job.Get("Target"), job.Get("Function"), job.Get("Arguments"))

	return nil
}
