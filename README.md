# AQA++
`AQA++` is an **unofficial** implementation of the [AQA psuedocode specification](https://filestore.aqa.org.uk/resources/computing/AQA-8520-TG-PC.PDF) in Go. It also supports some features that aren't in the spec (such as maps) in order to make it slightly more usable, hence the `++` in the name. It is interpreted and the design of the interpreter is based on the one described in the book [Writing An Interpreter in Go](https://interpreterbook.com/).

![Example program, calculates the GCD of two numbers](_media/gcd.png)

- [AQA++](#aqa)
  - [Demo](#demo)
  - [Syntax](#syntax)
    - [Specification](#specification)
    - [Additions](#additions)
  - [Bugs/Todo](#bugstodo)

## Demo
For now, if you want to try it out, you can go to [https://aqa.ollybr.repl.run/](https://aqa.ollybr.repl.run/) for a REPL and [https://repl.it/@ollybr/AQA](https://repl.it/@ollybr/AQA) if you want to write a file. Both of those are just REPL.ITs that download the executable from my website and run it.

For more examples, see the [_examples folder](_examples/).

## Syntax
### Specification
Everything in the AQA specification. One big difference is that there is no support for special characters at the moment. In the spec, it uses unicode characters such as `←` for assignment and `≥` for greater than or equal. In this version, only the ascii equivalents are supported.

| Specification | Equivalent | Purpose                                   |
|---------------|------------|-------------------------------------------|
| `←`           | `<-`       | Assignment: `a <- 10`                     |
| `≥`           | `>=`       | Greater than or equal: `10 >= 20 # false` |
| `≤`           | `<=`       | Less than or equal: `10 <= 20 # true`     |
| `≠`           | `!=`       | Not equal to: `10 != 20 # true`           |

### Additions
Additions to the spec (hence to `++`)
* NO UPPERCASE REQUIREMENTS SO THINGS DON'T NEED TO BE SCREAMED
* Use of `0x123edf` syntax to define hexadecimal numbers
* Use of `0b100000` syntax to define binary numbers
* More builtin functions, such as SQRT and FLOOR.
* Bitshifts using `>>` and `<<`

Also, it **WILL* support the following (to be added)
* Maps: using the `{` syntax `}`
* Automatic type conversion: adding an integer to a string wont cause an error.
* `FN`: similar to a subroutine, but an expression. This means `FN`s will be able to be passed around as arguments. 

## Bugs/Todo
- [ ] More tests
- [ ] Better type conversion system
- [ ] Lack of support for unicode