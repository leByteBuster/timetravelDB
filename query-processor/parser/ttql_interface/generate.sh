#!/bin/sh

alias antlr4='java -Xmx500M -cp "./antlr4-4.11.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool'
antlr4 -Dlanguage=Go TTQL.g4 -package ttql_interface -listener