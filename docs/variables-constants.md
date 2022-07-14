# nicer-syntax Variables and Constants

## Declaration

```perl
constant ConstNumber is number 10
variable VariableNumber is number
```

Variables and constants are declared in the same way, with differing keywords.
Variables are mutable, constants are immutable.
When declaring a variable, a value is optional at declaration-time; on the other hand, constants require a value at declaration time.

## Assignment

```perl
variable X is string
X is "Hello, world!"
```

The language uses `is` as its assignment operator.
Assignment without declaration is only allowed for variables.
Constants must be assigned to at declaration time.
