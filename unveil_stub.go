//go:build !openbsd
// +build !openbsd

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

func Unveil(path string, perms string) error {
	return nil
}

func UnveilBlock() error {
	return nil
}

func UnveilPaths(paths []string, perms string) error {
	return nil
}
