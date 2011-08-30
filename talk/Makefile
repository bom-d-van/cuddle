# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=makeslide
GOFILES=\
	makeslide.goo\

index.html: cuddle.slide makeslide slide.tmpl cuddle/snippets cuddle/hello cuddle/hellogae ../cuddle/http.go ../cuddle/name.go ../cuddle/db.go ../tmpl/root.html cuddle/cache cuddle/getroom
	./makeslide cuddle.slide > index.html

include $(GOROOT)/src/Make.cmd
