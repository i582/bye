The main feature of **bye** is the generation of a two-column representation for the code, where the code is on the
right side and the comment is on the right side.

To implement this, some features have been added to page formatting.

### Example

Let's take a look at an example:

```go (example.go)
// This is a simple golang program.
package main

// In order to *display* something on the screen, we need to include the `"fmt"` package.
import (
    "fmt"
)

func main() {
    // To display the string, use the special function `Println`.
    fmt.Println("Hello bye!")
}
```

The comments above were originally written in code, but after generation they were placed on the side of the code.

### Code comments

Some comments can be very long and so that they do not interfere with reading the code, they can be folded with an arrow
next to the comment.

You might want to hide **all comments** in order to *fully enjoy* the code. To do this, there is a `"toggle comment"`
button above the code on the right, which will hide all comments. There is also a button for copying a file next to it.

In some cases, you will want to leave a comment in the code, in which case you need to put a exclamation mark after the
symbols of the beginning of the comments

```php (example.php)
//! Use an exclamation mark to comment on your code.
/**
 * Or you can use multi-line comments.
 */
$a = 100;
echo $a;
```

### Themes

**bye** tries to be comfortable for customizing both *layout* and *styles*. By default, you already have two styles. The
current style is *standard light*, but you can switch to *dark* as well. To do this, change the value of the `theme`
field in the config.

```yaml (bye.yml)
theme: "dark",
```

### Tags

**bye** has the ability to give each page specific tags. Tags are set in the page description in the config in
the `tags` field. For each tag, its own page will be created, where all pages with this tag will be.

```yaml (bye.yml)
pages: [
    {
      name: "some name",
      title: "Some title",
# Here we set the required tags.
      tags: [
          "tag1",
          "tag2",
      ],
    },
],
```

### Footer

Let's take a look at the bottom of the page. There you will see a footer that contains a link to the next page (and the
previous one, if it exists). You can use the *left-right arrows* to navigate through the pages.