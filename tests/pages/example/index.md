[Source](https://doc.rust-lang.org/rust-by-example/hello.html).

This is the source code of the traditional Hello World program.

```rust (hello.rs)
//! This is a comment, and is ignored by the compiler

// This is the main function
//
// // Statements here are executed when the compiled binary is called
fn main() {
    // Print text to the console
    println!("Hello World!");
}
```

`println!` is a [*macro*](https://doc.rust-lang.org/rust-by-example/macros.html) that prints text to the console.

```bash
# A binary can be generated using the Rust compiler: `rustc`.
$ rustc hello.rs

#`rustc` will produce a `hello` binary that can be executed.
$ ./hello
Hello World!
```

