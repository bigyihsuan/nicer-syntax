# Language Structures

## Variables

### Variable Declaration

```perl
constant ConstNumber is number 10
variable VariableNumber is number
```

Variables and constants are declared in the same way, with differing keywords.
Variables are mutable, constants are immutable.
When declaring a variable, a value is optional at declaration-time; on the other hand, constants require a value at declaration time.

### Variable Assignment

```perl
variable X is string
X is "Hello, world!"
```

The language uses `is` as its assignment operator.

## Loops

### For Loops

```perl
for number I from 0 to 20, loop
    # do stuff
done
```

### While Loops

```perl
while Condition, loop
    # do stuff
done
```

## Ranges

Ranges are a big part of the language. They are *always* inclusive, no exceptions.

Ranges must contain a starting index and an ending index.
These indexes can be substituted by keywords `start` and `end`.
You can explicitly say which collections to base `start` and `end` by using `start of CollectionName` and `end of CollectionName`.

There is also an optional `every` clause, which will take every `N`-th element, starting from the starting index.
However, this is merely syntactic sugar, and a range like `from 10 to 20` is the same as `every 1-th from 10 to 20`.

```perl
from start to end
every 1-th from start to end
every 2-th from 10 to 20

constant List is list of numbers containing 0, 1, 1, 2, 3, 5, 8, and 13
for number I from start of List to end of List, do
    # do something
done
```

## If-Then-Else Statements

If statements are the main control flow.
There are three main forms: a single if, a single if-else, and a chained if-else.

```perl
# single if
if Condition, then
    do Something
done
# single if-else
if Condition1, then
    do Something
else, then
    do SomethingElse
done
# chained if-else
if Condition1, then
    do Something1
elseif Condition2, then
    do Something2
elseif Condition3, then
    do Something3
else, then
    do SomethingElse
done
```
