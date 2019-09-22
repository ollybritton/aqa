# Design
This document is just where I put temporary notes.

## Differences to Monkey
This is a list of differences/ideal changes to the Monkey programming language, the one defined in [Writing an Interpreter in Go](https://interpreterbook.com/).

* Unicode expressions "≠", "←" (though these will have ASCII versions)
* Constants (constant a <- 10)
* Infix variable statements (a <- 10 vs a = 10)
* Floats
* Type conversion
* OUTPUT being a statement vs a function (like Python 2 vs Python 3)
* USERINPUT being a statement
* Range using "1 TO 4" syntax

## Extensions to the Spec
This is a list of ideal things that would be useful to have in the language, that aren't defined in the specification.

* Maps, using the { construct }
* ~~Subroutines being expressions (so they can be passed around as first class citizens)~~
* Subroutines being implemented as they are in the spec, add a new FUNCTION construct which are like subroutines but expressions
* NO UPPERCASE REQUIREMENTS SO THINGS DON'T NEED TO BE SCREAMED
