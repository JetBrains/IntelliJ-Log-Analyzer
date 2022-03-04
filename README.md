# README

## About

Logs highlighter and analyzer. 
Program receives logs folder as an input and renders it for better usability.

## What can be parsed

For now, idea.log and ThreadDumps folders are being parsed. 
All unknown files are ignored. 
For adding your own file/folder type, extend the functionality as described in [Extending functionality](#Extending functionality) 

## Extending functionality

There are two types of information could be rendered in UI: *Static Info* and *Dynamic Info*. 
The only difference between them is timestamp existance.
For example, troubleshooting.txt is *Static Info*. idea.log, threadDumps folders, and jbr_err_in_PID are *Dynamic Info*, as all of them have timestamp assigned.

Your own filetype to analyze could be added as *Static Info* provider, or *Dynamic Info* provider. 
In both cases, you need to create new file in /backend/analyzer/Entities and 

