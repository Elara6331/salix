<img src="assets/icons/logo.svg" alt="Salix Logo" width="200">

&nbsp;

[![Go Reference](https://pkg.go.dev/badge/go.elara.ws/salix.svg)](https://pkg.go.dev/go.elara.ws/salix)
[![Go Report Card](https://goreportcard.com/badge/go.elara.ws/salix)](https://goreportcard.com/report/go.elara.ws/salix)

Salix (pronounced *say-lix*) is a Go templating engine inspired by [Leaf](https://github.com/vapor/leaf).

Salix's syntax is similar to Leaf and (in my opinion at least), it's much more fun to write than the Go template syntax. If you like this project, please star the repo. I hope you enjoy! :)

## Table of contents

- [Examples](#examples)
  - [Template](#template)
  - [API Usage](#api-usage)
- [Tags](#tags)
  - [Creating custom tags](#creating-custom-tags)
  - [`for` tag](#for-tag)
  - [`if` tag](#if-tag)
  - [`include` tag](#include-tag)
    - [Using the `include` tag with extra arguments](#using-the-include-tag-with-extra-arguments)
  - [`macro` tag](#macro-tag)
    - [Using the `macro` tag with extra arguments](#using-the-macro-tag-with-extra-arguments)
- [Functions](#functions)
  - [Global Functions](#global-functions)
  - [Adding Custom Functions](#adding-custom-functions)
- [Expressions](#expressions)
  - [Ternary Expressions](#ternary-expressions)
  - [Coalescing operator](#coalescing-operator)
  - [The `in` operator](#the-in-operator)
- [Markdown](#markdown)
- [Acknowledgements](#acknowledgements)

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

## Tags

In Salix, tags have full control over the Abstract Syntax Tree (AST), which allows them to do things the language wouldn't ordinarily allow. Salix's if statements, for loops, macros, includes, etc. are implemented as tags.

### Creating custom tags

You can extend the capabilities of Salix by creating custom tags. To create a custom tag, you need to implement the `salix.Tag` interface and add it to the tag map of your template or namespace using the `WithTagMap` method.

Salix tags follow a distinctive syntax pattern. They start with a pound sign (`#`), followed by a name and optional arguments. Tags can also enclose a block of content, and if they do, the block is terminated by an end tag with the same name. Here's an example of a template with both a `macro` tag and an `include` tag:

```
#macro("example"):
    Content
#!macro

#include("template.html")
```

In this example:

- The `macro` tag has a block, indicated by the content enclosed between `#macro("example"):` and `#!macro`.
- The `include` tag doesn't have a block; it simply includes the content of `template.html`.

### `for` tag

Salix's `for` tag is used for iterating over slices, arrays, and maps. It can assign one or two variables depending on your needs. When using a single variable, it sets that variable to the current element in the case of slices or arrays, or the current value for maps. With two variables, it assigns the first to the index (in the case of slices or arrays) or the key (for maps), and the second to the element or value, respectively. Here's an example of the for tag in action:

```
#for(id, name in users):
	Name: #(name)
	ID:   #(id)
#!for
```

### `if` tag

The `if` tag in Salix allows you to create conditional statements within your templates. It evaluates specified conditions and includes the enclosed content only if the condition is true. Here's an example:
```
#if(weather.Temp > 30):
	<p>It's a hot day!</p>
#elif(weather.Temp < 0):
	<p>It's freezing!</p>
#else:
	<p>The temperature is between 0 and 30</p<
#!if
```

### `include` tag

The include tag allows you to import content from other templates in the namespace, into your current template, making it easier to manage complex templates. Here's an example of the `include` tag:

```
#include("header.html")
```

#### Using the `include` tag with extra arguments

The `include` tag can accept extra local variables as arguments. Here's an example with a `title` variable:

```
#include("header.html", title = "Home")
```

These local variables will then be defined in the included template.

### `macro` tag

The macro tag is a powerful feature that allows you to define reusable template sections. These sections can be included later in the current template or in other templates that were included by the `include` tag. Here's an example of the macro tag:

```
#macro("content"): <!-- This defines a macro called "content" -->
    Content
#!macro

#macro("content") <!-- This inserts the content macro -->
```

When a macro tag has a block, it sets the macro's content. When it doesn't, it inserts the contents of the macro. In the above example, a macro is defined and then inserted.

#### Using the `macro` tag with extra arguments

Similar to the `include` tag, the `macro` tag can accept extra local variables as arguments. You can define these variables when including the macro. Here's an example:

```
#macro("content", x = 1, y = x + 2)
```

## Functions

Functions used in a template can accept any number of arguments but are limited to returning a maximum of two values. When a function returns two values, the second one must be an error value.

### Global Functions

Salix includes several useful global functions in all templates:

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

### Adding Custom Functions

You can include custom functions as variables using the WithVarMap method on templates or namespaces. Methods that fit the above conditions can also be used as template functions.

## Expressions

Salix's expressions mostly work like Go's, but there are some extra features worth mentioning.

### Ternary Expressions

Salix supports ternary expressions, which allow you to choose a value based on whether a condition is true. For example:

```
#(len(matches) > 1 ? "several matches" : "one match")
```

This example returns `"several matches"` if the length of matches is greater than one. Otherwise, it returns `"one match"`.

### Coalescing operator

The coalescing operator allows you to return a default value if a variable isn't defined. Here's an example:

```
<title>#(title | "Home")</title>
```

In this case, the expression will return the content of the `title` variable if it's defined. If not, it will return `"Home"` as the default value.

### The `in` operator

Salix's `in` operator allows you to check if a slice or array contains a specific element, or if a string contains a substring. Here's one example:

```
#("H" in "Hello") <!-- Returns true -->
```

## Markdown

Salix doesn't include a markdown rendering tag because I didn't want any non-stdlib dependencies. Instead, there's an implementation of a markdown salix tag using [goldmark](https://github.com/yuin/goldmark) at [go.elara.ws/salixmd](https://pkg.go.dev/go.elara.ws/salixmd).

## Acknowledgements

- [Pigeon](https://github.com/mna/pigeon): Salix uses a [PEG](https://en.wikipedia.org/wiki/Parsing_expression_grammar) parser generated by pigeon. Salix would've been a lot more difficult to write without it.
- [Leaf](https://github.com/vapor/leaf): Leaf was the first templaing language I ever used, and it inspired a lot of the syntax I've implemented in Salix because I found it really fun to use.

