[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/nalanj/confl)

# confl

Confl is a simple description language. It's related to languages like YAML,
JSON, and TOML. The focus with Confl is to keep the syntax as simple as
possible while also achieving a high amount of expressiveness.

## Example Document

Here's an example that shows most of the features of confl.

```
# Simple wifi configuration
device(wifi0)={
  network="Pretty fly for a wifi"
  key="Some long wpa key"
  dhcp=true

  dns=["10.0.0.1" "10.0.0.2"]
  gateway="10.0.0.1"

  vpn={host="12.12.12.12" user=frank pass=secret key=path(/etc/vpn.key)}
}
```

## The Tokens

### Numbers

Numbers are a series of digits, possibly including a decimal place. Hexidecimal
numbers are supported, prefixed with `0x` or `0X`.

```
12
12.5
0x12
```

### Strings

A string begins and ends with single or double quotes. Strings can contain
escaped quotation marks prefixed by `\`, and a `\` can be escaped by another
`\`. Strings ignore newlines, so long strings can span several lines:

```
"This is a string"
'This is a string'

"This is \"my\" string"
'This is \'your\' string'

"
  This
  is
  a
  long
  string
"
```

### Words

A word is a series of characters that does not start with a number and that
does not contain any spaces or quotation marks. Words are effectively short
strings that aren't surrounded by quotes, and are useful as identifiers.

Words cannot contain `"`, `'`, `=`, `,`, or any sort of whitespace.

```
word
a_word
```

Confl doesn't include an explicit boolean type because words can represent
booleans:

```
true
false
TRUE
FALSE
yes
no
```

### Maps

Maps are unordered key value pairs. The keys are always a string or a word.
The keys and values are separated by an equal `=` sign. Confl doens't care
about whitespace outside of strings, so a map may be all on a single line or
spread across multiple lines. Maps are surrounded by curly `{}` braces.

Map keys must be words or strings.

```
map={nested=map}

another_map = {
  another=map
}
```

All Confl documents are maps at the document level, so document level maps
exclude the curly braces.

### Lists

A list is an ordered series of values. Lists are surrounded by square brackets.
Lists are space delimited.

```
[this is "a" list of 7 "items"]
```

### Comments

Comments in Confl begin with the pound sign `#` and continue for the remainder
of the line.

```
# This is a comment

# Multiline comments are just
# multiple lines beginning with a pound
# sign
```

### Decorators

Decorators decorate other types to help them communicate more complex ideas. A
decorator is communicated as a word with a pair of attached parenthesis.

```
decorator(12)
```

A decorator can contain any other type, so long as the type would be valid in
that context without a decorator as well.

## Errors

Confl tries to do a good job with showing errors. The `Error()` function for a
`ParseError` simply returns the error message for the error, but there is an
additional `ErrorWithCode` function that includes information about the line
and location of the error. For example:

```
Illegal closing token: got }, expected EOF
Line 1: test=23 "also"=this}
                           ^
```

