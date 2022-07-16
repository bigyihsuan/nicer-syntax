# nicer-syntax Types

## Primitive Types

The language is strongly-typed, and explicit conversions are required for turning one type into another.
The language contains the following primitive types:

* `number`
* `boolean`
* `string`

`number` represents any real number.
`boolean` is either `true` or `false`.
`string` represents a string of 0 or more UTF8-encoded characters.

## Lists

Lists in the language are homogeneous lists of elements that can grow and shrink at will.
Lists are inherently generic: they do not care what is contained within. More information on generics below.
List literals are a comma-separated list of elements after the type declaration, with the last element preceded by `and`.

```perl
variable MemeNumbers is list of number containing 69, 420, and 9001, done
variable RangedList is list of number containing every 1-th from 10 to 20, done
variable OneElement is list of number containing 1, done
```

## Maps

Maps are collections mapping a unique key to a unique value.
Any type can be a key; any type can be a value.
Map literals are comma-separated lists of pairs after the type declaration, with each pair written as `key as value`.

Assigning a value to a key requires indexing the map.
If a key doesn't exist, it will be created automatically with a default value.

```perl
variable HouseNumbers is map of number to string containing 123 as "bob", 345 as "pat", and 420 as "dog", done
999-th of HouseNumbers is "rich" # map now contains key 999 with value "rich"
23-th of HouseNumbers is nothing # use default value of string (empty string) to init key 23
```

## Indexing and Ranging Lists, Strings, and Maps

Lists and strings can be indexed via a numeric index (`N-th from Collection`), as well as via a range (`every X-th from S to E from Collection`).
Indexing is 0-based, so `start` is equivalent to `0-th`.

Maps only support `from`-indexing through their key.

```perl
variable Hello is string "Hello World!"
start from Hello       # "H"
5-th from Hello        # " "
6-th to end from Hello # "World!"

variable Numbers is list of numbers containing every 1-th from 0 to 100, done
50-th from Numbers                      # 50
every 2-th from 1-th to end from Values # [1, 3, 5, ... 99]

variable HouseNumbers is map of string to number containing "bob" as 123, "pat" as 345, "dog" as 420, done
"bob"-th from HouseNumbers # 123
"pat"-th from HouseNumbers # 345
"abc"-th from HouseNumbers # 0 ; this is because the key does not exist and gives the default value of `number`
```

## Structs

Structs are data types composed of multiple fields.
These fields can be of any data type, including other structs.
Structs are a reference type, in that a variable containing a struct is a pointer to the actual struct somewhere else in memory.
Fields are completely optional, and can be wholy replaced with `nothing`.
Fields are separated by commas.

Structs can contain methods, as in functions that are associated with a struct and can be called by only that struct.
The are delimited from the field block by the `and can do` series of keywords.
Methods are completely optional, as seein the plain-old-data example below.

Accessing fields and methods uses the `FieldName of` operator before the struct's name.
A special keyword, `this`, refers to the instance of that struct when a method is called.

```perl
type BankAccount is struct containing
  variable AccountNumber is number,
  variable Balance is number,
  variable AccountHolder is string,
  variable Transactions is list of number,
and can do
  function Deposit is function, taking number Amount, doing
    Balance of this is Balance of this + Amount
  done
  function Withdraw is function, taking number Amount, returning boolean, doing
    if Balance of this <= Amount, do
      return false
    else do
      Balance of this is Balance of this - Amount
      return true
    done
  done
done
```

Plain-Old-Data (POD) structs are also possible:

```perl
type Tuple is struct of A and B containing
  constant first is A,
  constant second is B
done
```

Structs only containing methods are also possible:

```perl
type MethodsOnly is struct containing
    nothing
and can do
    function Something is function, doing
        Printline "Hello World!"
    done
done
```

Empty struct:

```perl
type Empty is struct containing
    nothing
done
```

Instantiating structs is done in two ways:

* by assigning to each field
* through a structure literal

```perl
type Tuple is struct of A, and B containing
  constant First is A,
  constant Second is B
done

# assignment to each field
variable Tup1 is Tuple of number, and string
First of Tup1 is 123
Second of Tup1 is "hello"

# structure literal
variable Tup1 is Tuple of number, and string containing
  First is 123
  Second is "hello"
done
```

## Default Values

Each of the types has a default value:

Type|Default
-|-
`number`|`0`
`string`|`""`
`list`|`nothing`
`map`|`nothing`
`struct`|`nothing`

`nothing` changes meaning based on what type it is exactly:

* A `list of E with nothing` is an empty list.
* A `map of K to V with nothing` is an empty map.
* A `struct` is a null reference.

## Generics

Lists and structs can take generic type parameters when they are declared.
For lists, this is a requirement.
The generic type is denoted by `of T` after the type name, where `T` is the name of a type.
Generics involving multiple different types can have a comma-separated list of type names after the first, with the last type being preceded by `and`.

```perl
type Numbers is list of number

type Container is struct of E containing
  # ...
done
```

## User-Defined Types

Using the `type` keyword, one can create customly-named types.
`type` is already used when creating structs.

```perl
type NumberList is list of number
type AddressBook is map of string to string

type Student is struct containing
    variable Name is string
    variable Id is number
    variable Classes is list of string
done
type Classroom is list of Student
```
