/*
 * Copyright (c)      2024  Runxi Yu <me@runxiyu.org>
 * Copyright (c) 2022-2024  Anirudh Oppiliappan <x@icyphox.sh>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
 * OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
 * SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package routes

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~runxiyu/legitrx/git"
)

func isGoModule(gr *git.GitRepo) bool {
	_, err := gr.FileContent("go.mod")
	return err == nil
}

func getDescription(path string) (desc string) {
	db, err := os.ReadFile(filepath.Join(path, "description"))
	if err == nil {
		desc = string(db)
	} else {
		desc = ""
	}
	return
}

func (d *deps) isIgnored(name string) bool {
	for _, i := range d.c.Repo.Ignore {
		if name == i {
			return true
		}
	}

	return false
}

type repoInfo struct {
	Git      *git.GitRepo
	Path     string
	Category string
}

func (d *deps) getAllRepos() ([]repoInfo, error) {
	repos := []repoInfo{}
	max := strings.Count(d.c.Repo.ScanPath, string(os.PathSeparator)) + 2

	err := filepath.WalkDir(d.c.Repo.ScanPath, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if de.IsDir() {
			// Check if we've exceeded our recursion depth
			if strings.Count(path, string(os.PathSeparator)) > max {
				return fs.SkipDir
			}

			if d.isIgnored(path) {
				return fs.SkipDir
			}

			// A bare repo should always have at least a HEAD file, if it
			// doesn't we can continue recursing
			if _, err := os.Lstat(filepath.Join(path, "HEAD")); err == nil {
				repo, err := git.Open(path, "")
				if err != nil {
					log.Println(err)
				} else {
					relpath, _ := filepath.Rel(d.c.Repo.ScanPath, path)
					repos = append(repos, repoInfo{
						Git:      repo,
						Path:     relpath,
						Category: d.category(path),
					})
					// Since we found a Git repo, we don't want to recurse
					// further
					return fs.SkipDir
				}
			}
		}
		return nil
	})

	return repos, err
}

func (d *deps) category(path string) string {
	return strings.TrimPrefix(filepath.Dir(strings.TrimPrefix(path, d.c.Repo.ScanPath)), string(os.PathSeparator))
}
