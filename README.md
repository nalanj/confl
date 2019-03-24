# confl

Confl is a simple configuration language. It's related to things like YAML,
JSON, and TOML.

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

```
{key=value "another key"="another value" number=12 map={another=map}}

{
  key=value
  "another key"="another value"
  number=12
  map = {
    another=map
  }
}
```

All Confl documents are maps at the document level, so the curly braces can be
excluded in that case.

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

A decorator can contain any other type.

Lists in a decorator can include or exclude the surrounding brackets.

```
list_decorator(test, 12)
```

Maps in a decorator can exclude the surrounded braces as well:

```
map_decorator(foo=bar blah=baz)
```
