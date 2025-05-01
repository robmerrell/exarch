# Exarch

Simple, feature incomplete semenatic searcher for Elixir.

## Description

Search elixir function calls, strings, and documentation. Mostly a toy program used
to play around with tree-sitter.

Function calls can either be searched with fully qualified names or partial matches.
Docs and Strings search with partial matches.

## Usage

```
NAME:
   exarch - Semantic Elixir Search

USAGE:
   exarch [global options] SEARCH_MODE SEARCH

DESCRIPTION:
   Searches recursively from the current directory.

   SEARCH_MODE can be one of:
   1. fncall - Search for function calls. This searches partial matches,
               and is able to handle aliases if given a fully qualified
               function. Eg. TestApp.Users.process can search for
               Users.process if TestApp.Users has been aliased.
   2. str - Search inside of strings.
   3. doc - search inside of documentation.

GLOBAL OPTIONS:
   --help, -h  show help
```
