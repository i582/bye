<!-- In PHP no specific way to create a *variable* explicitly, -->
<!-- as in `JavaScript` using the `let` keyword for example. -->

<?php

// To do this in PHP, you must assign a value to the variable.
$name = "John";
echo $name;
echo "\n";

// As with all languages, you can change the value of a variable.
$name = "Sasha";
echo $name;
echo "\n";

// PHP uses *dynamic typing*, which means that a variable can store any
// value with any type at any time. In *static typing*, for example, in C++, if a variable
// is declared as `int`, it can only store integers.
//
// Therefore, we can assign a different type of value to our variable,
// for example, a number.
$name = 1;
echo $name;
echo "\n";
