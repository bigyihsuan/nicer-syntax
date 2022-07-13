# nicer-syntax Operators

## Numeric Operators

Operator|Name|Notes
-|-|-
`+`|Addition|
`-`|Subtraction|Binary operator
`*`|Multiplication|
`/`|Division|
`%`|Modulo|
`^`|Exponentiation|
`-`|Negation|Unary operator

## Boolean Operators

Operator|Name|Notes
-|-|-
`and`|Logical AND|
`or`|Logical OR|
`not`|Logical NOT|Unary operator

## Comparison Operators

Operator|Name|Notes
-|-|-
`>`|Greater Than|
`>=`|Greater Than or Equal To|
`<`|Less Than|
`<=`|Less Than or Equal To|
`==`|Equal|
`!=`|Not Equal|

Comparison operators can be chained into a single expression if all signs follow the same direction. The following expressions are equivalent:

```perl
0 < N < Q
0 < N and N < Q # combinable
0 < N and Q > N # equivalent but not combinable
```

## Operator Precedence

The operators have the following precedence, from highest precedence to lowest:

1. `()`: Expression Grouping
2. `-`: Unary minus
3. `^`: Exponentiation
4. `* / %`: "Multiplication" operators
5. `+ -`: "Addition" operators
6. `> >= < <= == !=`: Comparison operators
7. `not`: Unary logical operators
8. `and or`: Binary logical operators

When multiple operators of the same precedence are chained together, they are evaluated left-to-right.

The following expression will be parsed as such, according to these precedence rules:

```perl
a + b * c > -d % e and not f - g / h ^ i + j == 0
# parsed as
((a + (b * c)) > ((-d) % e)) and (not (((f - ((g / (h ^ i))) + j) == 0))
```
