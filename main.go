/*
 * Copyright (c)      2024  Runxi Yu <me@runxiyu.org>
 * Copyright (c) 2022-2024  Anirudh Oppiliappan <x@icyphox.sh>
 *
 * SPDX-License-Identifier: AGPL-3.0-only
 * Overall, the work is licensed under AGPL-3.0-only. However, code by
 * Anirudh Oppiliappan is actually licensed under MIT/Expat. Check
 * git-blame(1).
 */

package main


import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"git.sr.ht/~runxiyu/legitrx/config"
	"git.sr.ht/~runxiyu/legitrx/routes"
)

func main() {
	var cfg string
	flag.StringVar(&cfg, "config", "./legitrx.yaml", "path to config file")
	flag.Parse()

	c, err := config.Read(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := UnveilPaths([]string{
		c.Dirs.Static,
		c.Repo.ScanPath,
		c.Dirs.Templates,
	},
		"r"); err != nil {
		log.Fatalf("unveil: %s", err)
	}

	mux := routes.Handlers(c)
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
