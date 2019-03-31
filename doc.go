/*
Package confl implements a simple description language. It's related to
languages like YAML, JSON, and TOML. The focus with Confl is to keep the
syntax as simple as possible while also achieving a high amount of
expressiveness.

## Example Document

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


It can be parsed using `confl.Parse()`:

	doc, err := confl.Parse(reader)

All confl documents are maps at their root.
*/
package confl
