/*
Package confl implements a simple description language. It's related to
languages like YAML, JSON, and TOML. The focus with Confl is to keep the
syntax as simple as possible while also achieving a high amount of
expressiveness.

Example Document

Here's an example that shows most of the features of confl.

	# Simple wifi configuration
	device(wifi0)={
		network="Pretty fly for a wifi"
		key="Some long wpa key"
		dhcp=true

		dns=["10.0.0.1" "10.0.0.2"]
		gateway="10.0.0.1"

		vpn={host="12.12.12.12" user=frank pass=secret key=path(/etc/vpn.key)}
	}

Parsing

Documents parsed using confl.Parse:

	doc, err := confl.Parse(reader)

Confl documents are always maps at their root.

Errors

Confl tries to do a good job with showing errors. The Error function for a
ParseError simply returns the error message for the error, but there is an
additional ErrorWithCode function that includes information about the line
and location of the error. For example:

	Illegal closing token: got }, expected EOF
	Line 1: test=23 "also"=this}

*/
package confl
