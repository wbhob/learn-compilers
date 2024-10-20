# number-sequence

This is a coding exercise to help me (and you) understand the basics of writing programming languages. The exercise is to write a program that parses programs written in a simple language and then executes them, returning the resulting string of numbers. In doing so, you will learn the basics of lexing, parsing, and executing programs. You'll build a basic compiler and interpreter that generates a list of tokens, parses them into an abstract syntax tree, and then executes the resulting instructions.

This exercise is inspired by a real-world use case to help a client create sequences of visual effects. In particular, it is used to represent a list of numbers that represent the number of frames to wait before displaying the next frame. By condensing a sequence of numbers into a compact syntax, it's easier for the client to create and edit sequences.

While this is a simple exercise, it is a good way to learn the basics of programming language design and implementation. You can use the foundational ideas to build more complex languages and compilers.

## Language Specification

The base case for the language is a sequence of numbers. For example:

```
2, 1, 1, 2, 2, 1
```

This is a sequence of 6 numbers. It represents a delay of 2 seconds, followed by a delay of 1 second, followed by a delay of 1 second, followed by 2 seconds, then 2 seconds, then 1 second.

The language allows you to condense this sequence into a more compact form using a few different operators. For example, you can use the `loop` operator to repeat a sequence multiple times. 

```
3x2, 1, 2x2, 1x2, 2
```

This is equivalent to the sequence `3, 3, 1, 2, 2, 1, 1, 2`.

You can also group numbers together using parentheses. 

```
(2, 1, 1), 2, 2, 1
```

This is equivalent to the sequence `2, 1, 1, 2, 2, 1`.

By combining these operators, you can create more complex sequences. 

```
(3x2, 1)x2, 1x2, 2
```

This is equivalent to the sequence `3, 3, 1, 3, 3, 1, 1, 1, 2`.

You can also nest parentheses and operators to create more complex sequences. 

```
((3x2, 1)x2, 1x2)x2, 2
```

This is equivalent to the sequence `3, 3, 1, 3, 3, 1, 1, 1, 3, 3, 1, 3, 3, 1, 1, 1, 2`.

It should also support decimals. 

```
1.5x2, 1
```

This is equivalent to the sequence `1.5, 1.5, 1`.

## Language Grammar

The language grammar is as follows:

```
<sequence> ::= <element> | <element> "," <sequence>
<element> ::= <number> | <group> | <loop>
<group> ::= "(" <sequence> ")"
<loop> ::= <element> "x" <integer>
<number> ::= <integer> | <decimal>
<integer> ::= <digit>+
<decimal> ::= <digit>+ "." <digit>+
<digit> ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
```

In laymans terms, this translates to:

```
<sequence> ::= either an **element** or an **element** followed by a **comma** and a **sequence**
<element> ::= is a **number**, a **group**, or a **loop**
<group> ::= a **sequence** wrapped in **parentheses**
<loop> ::= a **element** followed by an **x** and an **integer**
<number> ::= either an **integer** or a **decimal**
<integer> ::= one or more **digits**
<decimal> ::= one or more **digits**, followed by a **decimal point**, followed by one or more **digits**
<digit> ::= a single digit from **0 to 9**
```

## Language Tokens

The language tokens are the smallest units of the language. For example, the tokens for the sequence `3x2, 1` are `3`, `x`, `2`, `,`, and `1`.

The tokens are:

```
<token> ::= <number> | "x" | "," | "(" | ")"
```

That means when we start to build the list of tokens, we will look for a number, an `x`, a comma, a left parenthesis, or a right parenthesis. Every other character is either ignored or invalid.

## Language Lexer

The lexer is the first part of the compiler. It takes a string of text and converts it into a list of tokens. For example, the lexer for the sequence `3x2, 1` would return the list of tokens `["3", "x", "2", ",", "1"]`.

Now, your first thought may be to try to start to build the abstract syntax tree here. The syntax is simple enough, there are only a few rules to follow. And while you could write a program that does this, it would be hard to get right, and it would be difficult to maintain and extend.

Instead, we will build the lexer separately from the parser. The lexer will take the input string and convert it into a list of tokens, which will then be parsed into an abstract syntax tree. While this may seem overly simple and an unnecessary step, it makes it easier to build the parser later. It means that by the time the parser starts looking at the input string, we know that every character has been accounted for and is a meaningful part of the language. It also makes it easier to detect errors in the input string earlier, saving compute time and making it easier to understand where the error is.

A lexer can be as simple as separating an input string by whitespace. In this case, we may not have whitespace between tokens, so we need to take a step further. For this lexer, we will iterate over each character in the input string, and we will build up tokens from there. In this case, most of our tokens are single characters except for numbers, which may have digits before and after the decimal point. For this reason, we will need to keep track of the current number we are building up as we iterate over the input string, but we can just treat the decimal point a digit for the sake of constructing the token string.

Pseudocode for the lexer:

```
for each character in the input string:
    if the character is a digit:
        add it to the current number
    else if the character is a decimal point:
        add it to the current number
    else if the character is a comma OR x OR ( OR ):
        add the current character to the list of tokens
        reset the current number
    else if the character is whitespace:
        ignore it
    else:
        throw an error
```

The solution for the lexer is in the `lexer.go` file. Try implementing it yourself before looking at the solution.

## Language Parser

The parser is the next step in the compiler. It takes a list of tokens and converts it into an abstract syntax tree. For example, the parser for the sequence `3x2, 1` would return the abstract syntax tree along the lines of `loop{3, 2}, 1`.

The parser conforms the list of tokens to the grammar rules we defined earlier, building an abstract syntax tree that represents the input string. The abstract syntax tree is a tree of nodes that represent the structure of the program that can be passed to an interpreter or compiler to generate the resulting sequence of numbers.

The pseudocode for the parser:

```
for each token in the list of tokens:
    if the token is a number:
        add a value node to the tree with the number
    else if the token is a left parenthesis:
        find the corresponding right parenthesis
        build a tree for the sequence between the parentheses
        add a group node to the tree with the tree between the parentheses
        move cursor to the token after the right parenthesis
    else if the token is an x:
        add a loop node to the tree with sequence and count=next number
```

The solution for the parser is in the `parser.go` file. Try implementing it yourself before looking at the solution.

## Language Interpreter

The interpreter is the final step in the compiler. It takes an abstract syntax tree and converts it into a list of numbers. For example, the interpreter for the abstract syntax tree `loop{3, 2}, 1` would return the list of numbers `[3, 3, 2, 1]`.

The interpreter will recursively evaluate the abstract syntax tree, returning a list of numbers. The pseudocode for the interpreter:

```
for each node in the abstract syntax tree:
    if the node is a value node:
        add the value to the list of numbers
    else if the node is a group node:
        recurse on the sequence
    else if the node is a loop node:
        repeat the sequence the number of times specified by the loop
```

The solution for the interpreter is in the `interpreter.go` file. Try implementing it yourself before looking at the solution.

## Hints for Golang

- The `unicode` package has a function `IsSpace` that can be used to check if a character is a whitespace character.
- The `strconv` package has functions `Atoi` and `ParseFloat` that can be used to convert strings to integers and floats, respectively.

## Test cases

Tests are in the `main_test.go` file for golang. You can run them from the command line using `make test` or just `make`. You can get started by copying the `stubs.go` file to `main.go` and implementing the missing functions, as well as any functions you need. This foundational concept is portable to any language, so you can use AI to translate the code and tests into your favorite language.

## Acknowledgements

Thank you to [Caleb Miller](https://github.com/MilllerTime) for the inspiration and guidance in building this exercise, as well as providing the test cases.# learn-compilers
