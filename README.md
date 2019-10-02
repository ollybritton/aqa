# AQA++
![Example program, calculates the GCD of two numbers](_media/gcd.png)

`AQA++` is an **unofficial** implementation of the [AQA psuedocode specification](https://filestore.aqa.org.uk/resources/computing/AQA-8520-TG-PC.PDF) in Go. It also supports some features that aren't in the spec (such as maps) in order to make it slightly more usable, hence the `++` in the name. It is interpreted and the design of the interpreter is based on the one described in the book [Writing An Interpreter in Go](https://interpreterbook.com/).

It is interpreted diffentely to most production-ready languages, as there is no intermediate step of generating something like bytecode or instructions for a virtual machine.

## Demo
For now, if you want to try it out, you can go to [https://aqa.ollybr.repl.run/](https://aqa.ollybr.repl.run/) for a REPL and [https://repl.it/@ollybr/AQA](https://repl.it/@ollybr/AQA) if you want to write a file. Both of those are just REPL.ITs that download the executable from my website and run it.


## Bugs
See [BUGS](./BUGS.md)


## Design
See [DESIGN](./DESIGN.md)