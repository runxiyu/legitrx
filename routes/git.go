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
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
)

func (d *deps) InfoRefs(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	name = filepath.Clean(name)

	repo := filepath.Join(d.c.Repo.ScanPath, name)

	w.Header().Set("content-type", "application/x-git-upload-pack-advertisement")

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}

	billyfs := osfs.New(repo)
	loader := server.NewFilesystemLoader(billyfs)
	srv := server.NewServer(loader)
	session, err := srv.NewUploadPackSession(ep, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}

	ar, err := session.AdvertisedReferencesContext(r.Context())
	if errors.Is(err, transport.ErrRepositoryNotFound) {
		http.Error(w, err.Error(), 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ar.Prefix = [][]byte{
		[]byte("# service=git-upload-pack"),
		pktline.Flush,
	}

	if err = ar.Encode(w); err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}
}

func (d *deps) UploadPack(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	name = filepath.Clean(name)

	repo := filepath.Join(d.c.Repo.ScanPath, name)

	w.Header().Set("content-type", "application/x-git-upload-pack-result")

	upr := packp.NewUploadPackRequest()
	err := upr.Decode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Printf("git: %s", err)
		return
	}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}

	billyfs := osfs.New(repo)
	loader := server.NewFilesystemLoader(billyfs)
	svr := server.NewServer(loader)
	session, err := svr.NewUploadPackSession(ep, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}

	res, err := session.UploadPack(r.Context(), upr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}

	if err = res.Encode(w); err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: %s", err)
		return
	}
}
