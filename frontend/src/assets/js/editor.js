//showEditor dhows editor if it exists, or generate new editor if it does not exist
// @name is the id attribute for the editor
// @content is a async function which returns content to be displayed
function showEditor(name, content) {
    let id = getObjectID(name);
    let editors = $("#editors")
    $("#editors>div").hide()
    if (!$(`#${id}`).length && !$(`.${id}`).length) {
        if (content) cerateEditor()
    }
    editors.find(`.${id}`).show()
    return id

    async function cerateEditor() {
        editors.append(`    
            <div class=${id}>
                <div class="search-box" linked-editor="${id}"></div>
                <div id="${id}" class="editor">
                    <div class="loader">Loading...</div>
                </div>
            </div>
            `)
        const editor = ace.edit(id);
        editor.setOptions({
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
        editor.setValue(await content);
        editor.renderer.scrollToLine(Number.POSITIVE_INFINITY)
        editor.clearSelection();
        await createStyleGutterMarkers(0, editor.session.getLength())
        editor.session.foldAll(0, editor.session.getLength() - 4, 1);
        highlightEntriesTypes();
        editor.on("click", ThreadDumpLinkHandler)
        editor.on("click", IndexingDiagnosticLinkHandler)

        //Checks entryType of every line and highlight this line according to type.
        //Highlighting color is configured for every DynamicEntity on init()
        async function highlightEntriesTypes() {
            let mappedColors = JSON.parse(await window.go.main.App.GetEntityNamesWithLineHighlightingColors())
            let observer = new MutationObserver(function (e) {
                addHighlighting(e, mappedColors);
            });
            observer.observe($('.ace_gutter')[0], {subtree: true, childList: true});

            //addHighlighting is called on every mutation of the gutter, sets gutter's markers by createStyleMarkers(), and highlight the lines according to type from gutter
            async function addHighlighting(e, mappedColors) {
                e[0].target.childNodes.forEach(function (gutter) {
                    //todo: innerText is a hack, it is needed to get text line number from gutter position
                    let lineNumber = parseInt(gutter.innerText) - 1;
                    let lineAlreadyHighlighted = false;
                    Object.entries(mappedColors).forEach(([index]) => {
                        if (gutter.className.includes(getObjectID(index))) {
                            let marker = editor.session.getMarkers();
                            for (var i in marker) {
                                if (marker[i]["clazz"] === getObjectID(index) && marker[i]["range"]["start"]["row"] === lineNumber) {
                                    lineAlreadyHighlighted = true
                                    break;
                                }
                            }
                            if (!lineAlreadyHighlighted) {
                                editor.session.highlightLines(lineNumber, lineNumber, getObjectID(index), false)
                                lineAlreadyHighlighted = true
                            }
                        }

                    })
                })
            }
        }

        async function createStyleGutterMarkers(lineStart, lineEnd, mappedColors) {
            if (!mappedColors) {
                mappedColors = JSON.parse(await window.go.main.App.GetEntityNamesWithLineHighlightingColors())
            }
            let needle = /(^\s*<entryType>)(.*)(<\/entryType>\s*)/
            editor.$search.setOptions({
                needle: needle,
                caseSensitive: true,
                range: new ace.Range(lineStart, 0, lineEnd, Number.POSITIVE_INFINITY),
                wholeWord: false,
                regExp: true,
            });
            let range = editor.$search.findAll(editor.session)
            for (const rangeKey in range) {
                let groupName = editor.getSession().doc.getTextRange(range[rangeKey]).match(needle)[2]
                if (mappedColors[groupName]) {
                    if (mappedColors[groupName] !== true) {
                        let cssClass = getObjectID(groupName)
                        let cssContent = "position: absolute; opacity: 0.3; background-color:" + mappedColors[groupName] + ";"
                        addCssClass(cssClass, cssContent)
                        mappedColors[groupName] = true
                    }
                    editor.session.addGutterDecoration(range[rangeKey]["start"]["row"], getObjectID(groupName))
                }
                editor.session.replace(range[rangeKey], "")
            }
        }

        function addCssClass(className, content) {
            document.body.appendChild(
                Object.assign(
                    document.createElement("style"),
                    {textContent: ".ace_content ." + className + " {" + content + "}"})
            )
        }
    }
}