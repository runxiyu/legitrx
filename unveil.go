//go:build openbsd
// +build openbsd

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
	"golang.org/x/sys/unix"
	"log"
)

func Unveil(path string, perms string) error {
	log.Printf("unveil: \"%s\", %s", path, perms)
	return unix.Unveil(path, perms)
}

func UnveilBlock() error {
	log.Printf("unveil: block")
	return unix.UnveilBlock()
}

func UnveilPaths(paths []string, perms string) error {
	for _, path := range paths {
		if err := Unveil(path, perms); err != nil {
			return err
		}
	}
	return UnveilBlock()
}
