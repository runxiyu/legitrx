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
	"net/http"

	"git.sr.ht/~runxiyu/legitrx/config"
)

// Checks for gitprotocol-http(5) specific smells; if found, passes
// the request on to the git http service, else render the web frontend.
func (d *deps) Multiplex(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("rest")

	if r.URL.RawQuery == "service=git-receive-pack" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no pushing allowed!"))
		return
	}

	if path == "info/refs" &&
		r.URL.RawQuery == "service=git-upload-pack" &&
		r.Method == "GET" {
		d.InfoRefs(w, r)
	} else if path == "git-upload-pack" && r.Method == "POST" {
		d.UploadPack(w, r)
	} else if r.Method == "GET" {
		d.RepoIndex(w, r)
	}
}

func Handlers(c *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	d := deps{c}

	mux.HandleFunc("GET /", d.Index)
	mux.HandleFunc("GET /static/{file}", d.ServeStatic)
	mux.HandleFunc("GET /{name}", d.Multiplex)
	mux.HandleFunc("POST /{name}", d.Multiplex)
	mux.HandleFunc("GET /{name}/tree/{ref}/{rest...}", d.RepoTree)
	mux.HandleFunc("GET /{name}/blob/{ref}/{rest...}", d.FileContent)
	mux.HandleFunc("GET /{name}/log/{ref}", d.Log)
	mux.HandleFunc("GET /{name}/commit/{ref}", d.Diff)
	mux.HandleFunc("GET /{name}/refs/{$}", d.Refs)
	mux.HandleFunc("GET /{name}/{rest...}", d.Multiplex)
	mux.HandleFunc("POST /{name}/{rest...}", d.Multiplex)

	return mux
}
