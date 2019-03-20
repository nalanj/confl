# confl

Confl is a simple configuration language. It's related to things like YAML,
JSON, and TOML.

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
strings that don't contain any spaces. Words are useful as identifiers.

```
word
a_word
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

```
[this, is, "a", list, of, 7, "items"]
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
