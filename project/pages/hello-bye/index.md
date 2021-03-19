Hello. 👋

This is a *demo* page generated automatically.

If you haven't read the documentation yet, here's a [link](#) to it. It contains a *complete* description of the
capabilities.

Here we will *briefly* look at the capabilities of this tool.

## About.

**bye** was originally conceived as a *"by example"* documentation generator from code files.

However, it immediately became clear that it was not flexible. Therefore, **bye** has grown into a full page generator
from *markdown* files.

**bye** supports markdown format completely. You can also use markdown in the code comments, but there are limitations.

Documentation *"by example"* implies placing comments and code in *two columns*, which is, in general, the only
distinguishing feature from other static site generators. However, **bye** tries to take everything from this and make
it easier to view the code.

## Example

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

### Code comments setting

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

## Footer

Let's take a look at the bottom of the page. There you will see a footer that contains a link to the next page (if it
exists, and the previous one, if it exists). You can use the *left-right arrows* to navigate through the pages.

## Themes

**bye** tries to be comfortable for customizing both *layout* and *styles*. By default, you already have two styles. The
current style is *standard light*, but you can switch to *dark* as well. For this, there is a special switch in the
footer.

## Tags

**bye** has the ability to give each page specific tags. For each tag, an example page with that tag will be generated.
Tags are located above the title.

---

## After words

This is where the interesting moments associated with bye end. I hope this tool will help you in creating small sites
with examples for training.