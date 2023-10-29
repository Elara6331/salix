# Salix

Salix (pronounced *say-lix*) is a Go templating engine inspired by [Leaf](https://github.com/vapor/leaf).

Salix's syntax is similar to Leaf and (in my opinion at least), it's much more fun to write than the Go template syntax. If you like this project, please star the repo. I hope you enjoy! :)

## Examples

### Template

```html
<html>
    <head>
        <title>#(page.Title)</title>
    </head>
    <body>
        #for(i, user in users):
        <div>
            <h2>#(toLower(user.Name))</h2>
            <p>User ID: #(i)</p>
            #if(user.LoggedIn): <p>This user is logged in</p> #!if
            #if(user.IsAdmin): <p>This user is an admin!</p> #!if
            <p>Registered: #(user.RegisteredTime.Format("01-02-2006"))</p>
        </div>
        #!for
    </body>
</html>
```

### API Usage

```go
t, err := salix.New().ParseFile("example.salix.txt")
if err != nil {
  panic(err)
}

err = t.WithVarMap(vars).
    WithTagMap(tags).
    WithEscapeHTML(true).
    Execute(os.Stdout)
if err != nil {
  panic(err)
}
```

See the [examples](examples) directory for more examples.

## Functions

Functions that are used in a template can have any amount of arguments, but cannot have more than two return values. If a function has two return values, the second one must be an error value.

Salix includes the following default functions in all templates:

- `len(v any) int`: Returns the length of the value passed in. If the length cannot be found, it returns `-1`.
- `toUpper(s string) string`: Returns `s`, but with all characters mapped to their uppercase equivalents.
- `toLower(s string) string`: Returns `s`, but with all characters mapped to their lowercase equivalents.
- `hasPrefix(s, prefix string) bool`: Returns true if `s` starts with `prefix`.
- `trimPrefix(s, prefix string) string`: Returns `s`, but with `prefix` removed from the beginning.
- `hasSuffix(s, suffix string) bool`: Returns true if `s` ends with `suffix`.
- `trimSuffix(s, suffix string) string`: Returns `s`, but with `suffix` removed from the end.
- `trimSpace(s string) string`: Returns `s`, but with any whitespace characters removed from the beginning and end.
- `equalFold(s1, s2 string) bool`: Returns true if `s1` is equal to `s2` under Unicode case-folding, which is a general form of case-insensitivity.
- `count(s, substr string) int`: Returns the amount of times that `substr` appears in `s`.
- `split(s, sep string) []string`: Returns a slice containing all substrings separated by `sep`.
- `join(ss []string, sep string) string`: Returns a string with all substrings in `ss` joined by `sep`.

## What does the name mean?

Salix is the latin name for willow trees. I wanted to use a name related to plants since the syntax was highly inspired by Leaf, and I really liked the name Salix.

## Acknowledgements

- [Pigeon](https://github.com/mna/pigeon): Salix uses a [PEG](https://en.wikipedia.org/wiki/Parsing_expression_grammar) parser generated by pigeon. Salix would've been a lot more difficult to write without it.
- [Leaf](https://github.com/vapor/leaf): Leaf was the first templaing language I ever used, and it inspired a lot of the syntax I've implemented in Salix because I found it really fun to use.
