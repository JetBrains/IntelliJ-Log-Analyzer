IntelliJ Log Analyzer is an open source project created by me (Konstantin Annikov). 

Would you like to make it even better? That’s wonderful!

This page is created to help you start contributing. 

## Before you begin

Here is a short overview of technologies used in the project:

- GoLang is used as a backend and JavaScript (with JQuery) as a frontend.

- Backend and frontend are connected together via [Wails](https://github.com/wailsapp/wails) framework

- Logs are represented with the help of [Ace Editor](https://github.com/ajaxorg/ace) 

# Contributing

* Install GoLang, npm and wails command line utilities following [this guide](https://wails.io/docs/gettingstarted/installation) 
* Fork the repository and clone it to the local machine.
* Open the project with [GoLand](https://www.jetbrains.com/go/) or IntelliJ IDEA with Go plugin installed.

Here it is, you are all-set. Now you can run `wails dev` in project's root to start Log Analyzer in dev mode: 

![](https://i.imgur.com/jZu29uz.jpg)

# Where to start in the codebase

ALl the operations with log files (parsing, formatting, additional info generation) should be done in backend: `/backend/analyzer` 
Text highlighting, code folding, user-triggered events handling should be done in frontend: `/frontend/src/`  

## Extending the list of logs that could be parsed

There are two types of information that could be rendered in UI: *Static Info* and *Dynamic Info*.
The only difference between them is timestamp existance.
For example, troubleshooting.txt is *Static Info*. idea.log, threadDumps folders, and jbr_err_in_PID are *Dynamic Info*, as all of them have timestamp assigned.

### Adding your log file to the main logs viewer (Summary tool window)

There are 3 types of logs are displayed in the main editor on the below screenshot:

- build log (green highlighting)
- Indexing events (blue highlighting)
- Thread Dumps (red highlighting)

Each of the log types is added as a Dynamic Entity. All Dynamic Entities are combined and displayed in the main log viewer.

![](https://i.imgur.com/DuaQvKq.jpg)


If you would like to extend main log viewer with your type of logs, you'll need to add Dynamic Entity  via `CurrentAnalyzer.AddDynamicEntity` function on application initialization, for example: 

```go
    func init() {
        CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
            Name:           "Idea Log",     // Name of your log type to distinguish it from others.
            ConvertToLogs:  parseIdeaLog,   // function to convert idea.log file to analyzer.Logs struct
            CheckPath:      isIdeaLog,      // function that returns true in case a file/folder fits the requirement of your log type 
            GetDisplayName: getDisplayName, // functions that return string of how your log is represented in frontend 
        })
    }
```
It is recommended to create a file inside `backend/analyzer/entities` folder with the name of your entity. See [idea.log.go](backend/analyzer/entities/idea.log.go) as an example

### Adding info to the static info tool window 

Similar to *Dynamic Entities*, *Static Entities*  are combined and displayed in the Static info tool window: 

![](https://i.imgur.com/W9pJnDF.jpg)

To add your information to Static Info Tool window, it is needed to add Static Entiry via `` function on application initialization, for example: 

```go
func init() {
	CurrentAnalyzer.AddStaticEntity(analyzer.StaticEntity{
		Name:                "troubleshooting.txt",     // Name of your log type to distinguish it from others.
		ConvertToStaticInfo: parseTroubleshootingInfo,  // function to convert troubleshooting.txt file to analyzer.StaticInfo struct
		CheckPath:           isTroubleshootingInfo,     // function that returns true in case a file/folder fits the requirement of your log type
	})
}
```

## Extending the highlighting rules 

Highlighting rules are stored in [mode-idea_log.js](frontend/src/assets/js/lib/ace/mode-idea_log.js) file. Syntax and description of this file is available in [Defining Syntax Highlighting Rules](https://ace.c9.io/#nav=higlighter) section of Ace Editor documentation
Styles for highlighting rules are stored in [theme-light.css](frontend/src/assets/js/lib/ace/theme-light.css) file