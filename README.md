# IntelliJ Log Analyzer 
[![official JetBrains project](https://jb.gg/badges/official-flat-square.svg)](https://confluence.jetbrains.com/display/ALL/JetBrains+on+GitHub)

## About

Logs highlighter and analyzer for logs collected by **Help | Collect Logs and Diagnostic Data** action of any IntelliJ-based IDE.

Program receives logs folder as an input and renders it for better usability.

## How to use

- To tail the log of installed IDE, select it in the list of installed IDEs:

   <img src="https://i.imgur.com/IKYYEF3.png" width="500" alt="JetBrains Log Analyzer Select IDE">
- Drag&Drop file, folder, or archive to IntelliJ Log Analyzer window at any time to analyze.
  
  <img src="https://media.giphy.com/media/4LpM6HvPQ5mZs7pZTL/giphy.gif" width="500" alt="JetBrains Log Analyzer Select IDE">

- Click "Select directory" or "Select .zip" to open file/folder using OS file browser. 

## Demo 

[![JetBrains Log Analyzer](https://img.youtube.com/vi/BJf6XB2KGMg/0.jpg)](https://www.youtube.com/watch?v=BJf6XB2KGMg "JetBrains Log Analyzer")

## Builds

Download the latest release (Windows and macOS): https://github.com/annikovk/IntelliJ-Log-Analyzer/releases/latest/
Linux is not currently supported. 

## What can be parsed

For now,the following files are recognized: 
- idea.log files (including idea.log.1, etc)
- build-log folder
- threadDumps folders
- indexing-diagnostic folders

All unknown files are listed in **Other files** section.

License
=======
    Copyright 2022 Konstantin Annikov

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
