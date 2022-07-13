const styleMarkerNeedle = /(^\s*<entryType>)(.*)(<\/entryType>\s*)/
const editors = $("#editors");


//showEditor dhows editor if it exists, or generate new editor if it does not exist
// @name is the id attribute for the editor
// @content is a async function which returns content to be displayed
async function showEditor(name, content) {
    let id = getObjectID(name);
    let highlightingColors = JSON.parse(await window.go.main.App.GetEntityNamesWithLineHighlightingColors())
    let isHighlighted = false
    $("#editors>div").hide()
    if (!$(`#${id}`).length && !$(`.${id}`).length) {
        if (content) {
            let editor = await cerateEditor()
            setEditorOptions(editor, {
                foldOnLoad: true,
                fontSize: await window.go.main.App.GetSetting("EditorFontSize"),
                useSoftWrap: await window.go.main.App.GetSetting("EditorDefaultSoftWrapState")
            })
            bindEditorListeners(editor)
        }
    }
    editors.find(`.${id}`).show()
    return id

    function bindEditorListeners(editor) {
        window.runtime.EventsOn("LogsUpdated", function (s) {
            let insertedRange = editor.session.insert({
                row: editor.session.getLength(),
                column: 0
            }, s);
            let marker = editor.session.highlightLines(insertedRange.row, insertedRange.row, "justchangedline", false)
            setTimeout(() => {
                editor.session.removeMarker(marker.id)
            }, 2000)
            if (editor.renderer.layerConfig.lastRow >= editor.session.getLength() - 7) {
                editor.renderer.scrollToLine(Number.POSITIVE_INFINITY)
            }

        })
        editor.session.on("changeScrollTop", async function () {
        })
        editor.on("click", ThreadDumpLinkHandler)
        editor.on("click", IndexingDiagnosticLinkHandler)
        editor.renderer.on('afterRender', function () {
            createStyleGutterMarkers(editor, highlightingColors).unfolded
        })
    }

    async function setEditorOptions(editor, optList) {
        editor.setFontSize(optList.fontSize)
        editor.session.setUseWrapMode(optList.useSoftWrap)
        editor.session.foldAll(0, editor.session.getLength() - 2, 1)
        createStyleGutterMarkers(editor, highlightingColors).all
        isHighlighted = true;
        editor.session.foldAll(0, editor.session.getLength() - 2, 1)
        console.log("Folded all");
    }

    async function cerateEditor() {
        editors.append(`    
            <div class=${id}>
                <div class="search-box" linked-editor="${id}"></div>
                <div id="${id}" class="editor">
                    <div class="loader">Loading...</div>
                </div>
            </div>
            `)
        let editor = ace.edit(id);
        await editor.setOptions({
            mode: 'ace/mode/idea_log',
            theme: "ace/theme/idealog",
            readOnly: true,
            selectionStyle: "text",
            showLineNumbers: true,
            showGutter: true,
            showPrintMargin: false,
            highlightSelectedWord: true,
            scrollPastEnd: 0.05,
        })
        editor.execCommand('find');
        window.runtime.LogDebug("Fetching content for " + name)
        editor.setValue(await content);
        editor.renderer.scrollToLine(Number.POSITIVE_INFINITY)
        editor.clearSelection();
        console.log("Created editor for " + name)
        return editor
        //Checks entryType of every line and highlight this line according to type.
        //Highlighting color is configured for every DynamicEntity on init()
    }
    function createStyleGutterMarkers(e, mappedColors) {
        let editor = ace.edit(e)
        let foldWidgets = editor.session.foldWidgets
        let rowsOnScreen = editor.getLastVisibleRow() - editor.getFirstVisibleRow()
        function hghlightUnfoldedLines() {
            for (const i in foldWidgets) {
                if (foldWidgets[i] === "") {
                    highlightLine(editor, i, mappedColors)
                }
            }
        }
        function highlightAllLines() {
            for (let i = editor.session.getLength(); i >= 0; i--) {
                if (foldWidgets[i] === "" || foldWidgets[i] === "start") {
                    highlightLine(editor, i, mappedColors)
                }
            }
        }
        function highlightAllVisibleLines() {
            for (let i = editor.getFirstVisibleRow(); i < editor.getLastVisibleRow(); i++) {
                if (foldWidgets[i] === "") {
                    if (editor.session.getLine(i).includes("entryType")) {
                        highlightLine(editor, i, mappedColors)
                    }
                }
            }
        }
        function highlightLine(editor, lineNumber, mappedColors) {
            if (!editor.session.getLine(lineNumber).match(styleMarkerNeedle)) {
                return
            }
            let lineContent = editor.session.getLine(lineNumber).match(styleMarkerNeedle)[0]
            let groupName = lineContent.match(styleMarkerNeedle)[2]
            let startPosition = lineContent.indexOf(lineContent.match(styleMarkerNeedle)[0])
            let endPosition = lineContent.match(styleMarkerNeedle)[0].length
            let range = new ace.Range(lineNumber, startPosition, lineNumber, endPosition)
            if (mappedColors[groupName] !== true) {
                let cssClass = getObjectID(groupName)
                let cssContent = "position: absolute; opacity: 0.3; background-color:" + mappedColors[groupName] + ";"
                addCssClass(cssClass, cssContent)
                mappedColors[groupName] = true
            }
            if (mappedColors[groupName]) {
                editor.session.highlightLines(lineNumber, lineNumber, getObjectID(groupName), false)
                // editor.session.addGutterDecoration(lineNumber, getObjectID(groupName))
            }
            if (editor.session.foldWidgets[lineNumber] === "start") {
                // editor.session.replace(range, "")
                // editor.getSession().foldAll()
            }
            editor.session.replace(range, "")


        }
        function addCssClass(className, content) {
            document.body.appendChild(
                Object.assign(
                    document.createElement("style"),
                    {textContent: ".ace_content ." + className + " {" + content + "}"})
            )
        }

        return {
            all: highlightAllLines(),
            visible: highlightAllVisibleLines(),
            unfolded: hghlightUnfoldedLines()
        };
    }
}
