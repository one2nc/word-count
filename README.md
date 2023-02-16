# word-count

Write a program that implements Unix command wc (word count) like functionality. Read the Linux man page for wc in case you don't know what it does.

## Setup instructions

- To build a binary of the program use the `make build`.
- Install the word count program, using the `make install` command.
- Tests can be run using the `make test` command.

## Usage Instructions

The program can now be used in the terminal using the `gowc` command.

- It takes `-l` - lines, `-c`- characters and the `-w`- words flags in any order.
- If no flags are used counts of all line,word and char count will be displayed.
- Piping to the program is also possible e.g `cat README.md | gowc`
- You can also start the program and type in the terminal like so

```bash
    $ gowc
some randomtext here #press `CTRL + D` to exit.
       1       3      21 #result
```
