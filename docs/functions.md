# nicer-syntax Functions

## Function Declaration

Functions in the language are always anonymous.
A minimal, anonymous function that does nothing is as follows:

```perl
function doing
    nothing # an empty body can be replaced with keyword `nothing`
done
```

Naming a function is similar to creating a variable:

```perl
function DoesNothing is function doing
    nothing
done
```

## Function Parameters

Function parameters are a comma-separated list of type name-parameter name pairs, with the last pair being preceded by `and`.

```perl
function DoesSomething is function, taking number Foo, string Bar, and map of number and boolean Quux doing
    # function body here
done
```

Function parameters by default are mutable and can be used like variables;
you can force immutability by prepending `constant` before the type name:

```perl
function DoesSomethingConstant is function, taking constant number Foo, constant string Bar, and constant map of number and boolean Quux doing
    # function body here
done
```

## Calling Functions

Functions can be called by invoking its name, after the keyword `do`, then its parameters preceded with `to`.
If a function takes no arguments, it can omit the `to` clause, or give it `to nothing`.

```perl
function SaysHello is function, returning string, doing
    return "Hello!"
done

do PrintLine to do SaysHello to nothing, and do SaysHello + " World!" # PrintLine(SaysHello(), SaysHello() + " World!")
```

## Functions as Parameters

You can pass functions as a parameter to a function.

```perl
# typical filter function, takes a list and a function that returns boolean
function Filter is function, taking list of E Elements,
            and function, taking E, and returning boolean, FilteringFunction,
            returning list of E, doing
    variable Output is list of E containing nothing
    for number I from start of Elements to end of Elements, loop
        if do FilteringFunction to I-th of Elements, then
            # Append is a built-in that appends to the end of the collection
            do Append to Output and I-th of Elements
        done
    done
    return Output
done
```
