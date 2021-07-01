package main

import (
	"context"
	"log"
	"os"

	"gitlab-master.nvidia.com/itml-public/salt-api-go"
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

	c.Minions(ctx, func(id string, grains salt.Response) error {
		log.Printf("ID = %s, osfinger = %s", id, grains.Get("osfinger"))

		return nil
	})

	c.Logout(ctx)
}
