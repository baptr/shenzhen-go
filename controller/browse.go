// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/shenzhen-go/model"
	"github.com/google/shenzhen-go/view"
)

type dirBrowser struct{}

// DirBrowser serves a way of visually navigating the filesystem.
var DirBrowser dirBrowser

func (dirBrowser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s browse: %s", r.Method, r.URL)

	path := r.URL.Path
	_, reload := r.URL.Query()["reload"]
	if g, ok := loadedGraphs[path]; ok && !reload {
		Graph(g, w, r)
		return
	}

	base := filepath.Join(".", path)
	f, err := os.Open(base)
	if err != nil {
		log.Printf("Couldn't open: %v", err)
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Printf("Couldn't stat: %v", err)
		http.NotFound(w, r)
		return
	}
	if !fi.IsDir() {
		g, err := model.LoadJSON(f, base)
		if err != nil {
			log.Printf("Not a directory or a valid JSON-encoded graph: %v", err)
			http.ServeContent(w, r, f.Name(), fi.ModTime(), f)
			return
		}
		loadedGraphs[path] = g
		Graph(g, w, r)
		return
	}

	if nu := r.URL.Query().Get("new"); nu != "" {
		// Check for an existing file.
		nfp := filepath.Join(base, nu)
		if _, err := os.Stat(nfp); !os.IsNotExist(err) {
			log.Printf("Asked to create %q but it already exists", nfp)
			http.Error(w, "File already exists", http.StatusNotModified)
			return
		}
		path = filepath.Join(path, nu)
		pkgp, err := GuessPackagePath(nfp)
		if err != nil {
			log.Printf("Guessing a package path: %v", err)
		}
		loadedGraphs[path] = model.NewGraph(nfp, pkgp)
		http.Redirect(w, r, path+"?props", http.StatusSeeOther)
		return
	}

	fis, err := f.Readdir(0)
	if err != nil {
		log.Printf("Couldn't readdir: %s", err)
		http.NotFound(w, r)
		return
	}

	var e []view.DirectoryEntry
	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		e = append(e, view.DirectoryEntry{
			IsDir: fi.IsDir(),
			Name:  fi.Name(),
			Path:  filepath.Join(path, fi.Name()),
		})
	}

	view.Browse(w, base, e)
}
