// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// developer package defines a webserver for serving compiled code that automatically re-compiles
// on refresh.
package developer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// Run runs the development webserver on localhost at the given addr.
func Run(addr string, rootSourceFilePath string, debug bool, vcsDevelopmentDirectories []string) bool {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	var transaction *developTransaction = nil
	name := filepath.Base(rootSourceFilePath)

	serveRuntime := func(w http.ResponseWriter, r *http.Request) {
		transaction = newDevelopTransaction(rootSourceFilePath, vcsDevelopmentDirectories, addr, name)
		transaction.Start(w, r)
	}

	serveAndRecompile := func(w http.ResponseWriter, r *http.Request) {
		transaction.Build(w, r)
	}

	serveSourceMap := func(w http.ResponseWriter, r *http.Request) {
		transaction.ServeSourceMap(w, r)
	}

	serveSourceFile := func(w http.ResponseWriter, r *http.Request) {
		transaction.ServeSourceFile(w, r)
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/"+name+".js", serveRuntime).Methods("GET")
	rtr.HandleFunc("/"+name+".develop.js", serveAndRecompile).Methods("GET")
	rtr.HandleFunc("/"+name+".develop.js.map", serveSourceMap).Methods("GET")
	rtr.HandleFunc("/source/{path:.+}", serveSourceFile).Methods("GET")

	http.Handle("/", rtr)

	highlight := color.New(color.FgHiWhite, color.Underline).SprintFunc()
	fmt.Printf("Serving development server for project %v on %v at %v\n", highlight(rootSourceFilePath), addr, highlight("/"+name))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error running develop: %v", err)
		return false
	}

	return true
}