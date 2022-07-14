# nicer-syntax Loops and Ranges

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

## Loops

### For Loops

For loops iterate over a range or collection (list or map).
Looping over a range counts up from the lower bound to the upper bound, inclusive.
Looping over a collection gets a different element based on the collection:

* Lists get the element directly.
* Maps get the key.

```perl
# ranged for-loop
for number I from 0 to 5, loop
    do Printline to I # prints 0 1 2 3 4 5
done

# collection for-loop: lists
variable Strings is list of string containing "foo", "bar", and "baz"
for string S from start of Strings to end of Strings, loop
    do Printline to S # prints foo bar baz
done

variable Mapping is map of string to number containing "foo" as 123, "bar" as 456, "baz" as 789
for string Key from start to end of Mapping, loop
    do Printline to Key
    do Printline to Key-th from Mapping
    # prints foo 123 bar 456 baz 789
done
```

### While Loops

```perl
while Condition, loop
    # do stuff
done
```
