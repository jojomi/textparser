# textparser
When a regexp is just not convenient…

This library is meant as a quick way of parsing a string by code, but without implementing all the low-level handling from scratch.

## Goals

* be more dynamic than a regexp, easier to branch and with "memory" by allowing Go variables as storage
* support unicode strings throughout
* provide rune based as well as string based functions
* emit relevant and useful error messages

## Non-goals

* streaming support
* high-performance parsing

## Example

### Parse Key Value List

Input:

```
name: Tom Williams
role: CEO
os: Debian Linux
```

Parsing Code:

