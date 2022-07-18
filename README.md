# nicer-syntax

"when you're a tutor for cs and you dictate the code to the student but they just write your words down instead of code, the language" --- Me, bigyihsuan

## Goals

* Make the langauge's syntax as "nice" as possible (to me at least).
* Learn how to do anonymous functions, functions as values, etc.
* Learn how to do whitespace-indentation.
* Write an interpreter in ~~Python Racket. (for [Racket's Summer #lang Party][racket])~~ Golang to start.
* Once complete, port this interpreter to another language (maybe Racket).

## Random Thoughts

* Use many, many keywords
* Keywords should be full words (no abbreviations)
* Symbols should be used sparingly
* Explicit im/mutability
* Types
  * Number
  * Boolean
  * String
  * List
  * Map
  * Struct
* Type names come after the variable
* Variables must be explicitly typed, no exceptions
* Identifiers *must* be PascalCase. Tokens starting with lowercase letters are assumed to be keywords. (shamlessly stolen without much thought from [Cognate][cognate])

## Random Useful Links?

* <https://github.com/db47h/lex>

[racket]: https://racket.discourse.group/t/summer-lang-party/1128?u=spdegabrielle
[cognate]: https://cognate-lang.github.io/
