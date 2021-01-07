// Code generated by qtc from "basepage.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// This is a base page template. All the other template pages implement this interface.
//

//line basepage.qtpl:3
package main

//line basepage.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line basepage.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line basepage.qtpl:4
type page interface {
//line basepage.qtpl:4
	Title() string
//line basepage.qtpl:4
	StreamTitle(qw422016 *qt422016.Writer)
//line basepage.qtpl:4
	WriteTitle(qq422016 qtio422016.Writer)
//line basepage.qtpl:4
	Head() string
//line basepage.qtpl:4
	StreamHead(qw422016 *qt422016.Writer)
//line basepage.qtpl:4
	WriteHead(qq422016 qtio422016.Writer)
//line basepage.qtpl:4
	Body() string
//line basepage.qtpl:4
	StreamBody(qw422016 *qt422016.Writer)
//line basepage.qtpl:4
	WriteBody(qq422016 qtio422016.Writer)
//line basepage.qtpl:4
	Scripts() string
//line basepage.qtpl:4
	StreamScripts(qw422016 *qt422016.Writer)
//line basepage.qtpl:4
	WriteScripts(qq422016 qtio422016.Writer)
//line basepage.qtpl:4
}

// page prints a page implementing page interface.

//line basepage.qtpl:14
func streampageTemplate(qw422016 *qt422016.Writer, p page) {
//line basepage.qtpl:14
	qw422016.N().S(`
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
<meta name="description" content="Dreamvat Hub Host" />
<meta name="author" content="Lukiya Chen" />
<title>`)
//line basepage.qtpl:22
	p.StreamTitle(qw422016)
//line basepage.qtpl:22
	qw422016.N().S(`</title>
<link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.3.1/css/bootstrap.min.css" rel="stylesheet" />
<link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.11.2/css/all.min.css" rel="stylesheet" />
<link href="/css/site.min.css" rel="stylesheet">
`)
//line basepage.qtpl:26
	p.StreamHead(qw422016)
//line basepage.qtpl:26
	qw422016.N().S(`
</head>
<body class="text-center">
`)
//line basepage.qtpl:29
	p.StreamBody(qw422016)
//line basepage.qtpl:29
	qw422016.N().S(`
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.1/jquery.slim.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.3.1/js/bootstrap.min.js"></script>
`)
//line basepage.qtpl:32
	p.StreamScripts(qw422016)
//line basepage.qtpl:32
	qw422016.N().S(`
</body>
</html>
`)
//line basepage.qtpl:35
}

//line basepage.qtpl:35
func writepageTemplate(qq422016 qtio422016.Writer, p page) {
//line basepage.qtpl:35
	qw422016 := qt422016.AcquireWriter(qq422016)
//line basepage.qtpl:35
	streampageTemplate(qw422016, p)
//line basepage.qtpl:35
	qt422016.ReleaseWriter(qw422016)
//line basepage.qtpl:35
}

//line basepage.qtpl:35
func pageTemplate(p page) string {
//line basepage.qtpl:35
	qb422016 := qt422016.AcquireByteBuffer()
//line basepage.qtpl:35
	writepageTemplate(qb422016, p)
//line basepage.qtpl:35
	qs422016 := string(qb422016.B)
//line basepage.qtpl:35
	qt422016.ReleaseByteBuffer(qb422016)
//line basepage.qtpl:35
	return qs422016
//line basepage.qtpl:35
}

// Base page implementation. Other pages may inherit from it if they need
// overriding only certain page methods

//line basepage.qtpl:40
type basePage struct {
	Username string
}

//line basepage.qtpl:43
func (x *basePage) StreamTitle(qw422016 *qt422016.Writer) {
//line basepage.qtpl:43
	qw422016.N().S(`Dreamvat Passport`)
//line basepage.qtpl:43
}

//line basepage.qtpl:43
func (x *basePage) WriteTitle(qq422016 qtio422016.Writer) {
//line basepage.qtpl:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line basepage.qtpl:43
	x.StreamTitle(qw422016)
//line basepage.qtpl:43
	qt422016.ReleaseWriter(qw422016)
//line basepage.qtpl:43
}

//line basepage.qtpl:43
func (x *basePage) Title() string {
//line basepage.qtpl:43
	qb422016 := qt422016.AcquireByteBuffer()
//line basepage.qtpl:43
	x.WriteTitle(qb422016)
//line basepage.qtpl:43
	qs422016 := string(qb422016.B)
//line basepage.qtpl:43
	qt422016.ReleaseByteBuffer(qb422016)
//line basepage.qtpl:43
	return qs422016
//line basepage.qtpl:43
}

//line basepage.qtpl:44
func (x *basePage) StreamBody(qw422016 *qt422016.Writer) {
//line basepage.qtpl:44
}

//line basepage.qtpl:44
func (x *basePage) WriteBody(qq422016 qtio422016.Writer) {
//line basepage.qtpl:44
	qw422016 := qt422016.AcquireWriter(qq422016)
//line basepage.qtpl:44
	x.StreamBody(qw422016)
//line basepage.qtpl:44
	qt422016.ReleaseWriter(qw422016)
//line basepage.qtpl:44
}

//line basepage.qtpl:44
func (x *basePage) Body() string {
//line basepage.qtpl:44
	qb422016 := qt422016.AcquireByteBuffer()
//line basepage.qtpl:44
	x.WriteBody(qb422016)
//line basepage.qtpl:44
	qs422016 := string(qb422016.B)
//line basepage.qtpl:44
	qt422016.ReleaseByteBuffer(qb422016)
//line basepage.qtpl:44
	return qs422016
//line basepage.qtpl:44
}

//line basepage.qtpl:45
func (x *basePage) StreamHead(qw422016 *qt422016.Writer) {
//line basepage.qtpl:45
}

//line basepage.qtpl:45
func (x *basePage) WriteHead(qq422016 qtio422016.Writer) {
//line basepage.qtpl:45
	qw422016 := qt422016.AcquireWriter(qq422016)
//line basepage.qtpl:45
	x.StreamHead(qw422016)
//line basepage.qtpl:45
	qt422016.ReleaseWriter(qw422016)
//line basepage.qtpl:45
}

//line basepage.qtpl:45
func (x *basePage) Head() string {
//line basepage.qtpl:45
	qb422016 := qt422016.AcquireByteBuffer()
//line basepage.qtpl:45
	x.WriteHead(qb422016)
//line basepage.qtpl:45
	qs422016 := string(qb422016.B)
//line basepage.qtpl:45
	qt422016.ReleaseByteBuffer(qb422016)
//line basepage.qtpl:45
	return qs422016
//line basepage.qtpl:45
}

//line basepage.qtpl:46
func (x *basePage) StreamScripts(qw422016 *qt422016.Writer) {
//line basepage.qtpl:46
}

//line basepage.qtpl:46
func (x *basePage) WriteScripts(qq422016 qtio422016.Writer) {
//line basepage.qtpl:46
	qw422016 := qt422016.AcquireWriter(qq422016)
//line basepage.qtpl:46
	x.StreamScripts(qw422016)
//line basepage.qtpl:46
	qt422016.ReleaseWriter(qw422016)
//line basepage.qtpl:46
}

//line basepage.qtpl:46
func (x *basePage) Scripts() string {
//line basepage.qtpl:46
	qb422016 := qt422016.AcquireByteBuffer()
//line basepage.qtpl:46
	x.WriteScripts(qb422016)
//line basepage.qtpl:46
	qs422016 := string(qb422016.B)
//line basepage.qtpl:46
	qt422016.ReleaseByteBuffer(qb422016)
//line basepage.qtpl:46
	return qs422016
//line basepage.qtpl:46
}
