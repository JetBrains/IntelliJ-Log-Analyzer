const render = async () => {
    await clearToolWindows();
    await redrawEditors()
    $(".group-label>.folding-icon").click();
}
const redrawEditors = async () => {
    $("#editors>div").remove();
    await renderMainScreen();
    setSidebarState()
}
const clearToolWindows = async () => {
    let tabs = $("#toolWindows-buttons").find(".toolWindowButton")
    let toolWindows = $("#toolWindows").find(".toolWindow")
    toolWindows.remove();
    tabs.remove();
}

async function renderMainScreen() {
    await showToolWindow("Summary", "filters", "top", "Main Editor", window.go.main.App.GetSummary())
    if (await window.go.main.App.GetStaticInfo()) {
        await showToolWindow("Static Info", "staticinfo", "bot", "", window.go.main.App.GetStaticInfo())
    }
    window.mainEditorID = showEditor("Main Editor", window.go.main.App.GetLogs());
}

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
            theme: "ace/theme/light",
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
        createStyleGutterMarkers(0, editor.session.getLength())
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

// showToolWindow adds tool window on the left.
// @name is the Title of Tool Window
// @position is one of top/bot
// @fillFunction is a promise for Go function that returns a string of Tool Windiow content
async function showToolWindow(name, cssClass, position, linkedEditor, fillFunction) {
    let tabs = $("#toolWindows-buttons ." + position);
    let tabsElements = tabs.find(`.toolWindowButton`);
    let toolWindows = $("#toolWindows ." + position + " .toolWindow")
    let id = getObjectID(name);
    if (!getToolWindowTabElement()) {
        createToolWindowTabElement()
        await createToolWindowContent()
        return showToolWindow(name, cssClass, position, linkedEditor, fillFunction)
    }
    selectToolWindowTab()
    showToolWindowContent()

    function getToolWindowTabElement() {
        let a;
        tabsElements.each(function () {
            let attr = $(this).attr('target')
            if (attr === id) {
                a = $(this)
                return a
            }
        })
        return a
    }

    function selectToolWindowTab() {
        tabsElements.each(function () {
            if ($(this).attr('target') === id) {
                $(this).addClass("active")
            } else {
                $(this).removeClass("active")
            }
        })
    }

    function showToolWindowContent() {
        toolWindows.each(function () {
            if ($(this).prop('id') === id) {
                $(this).show()
                $(this).parent().show()
            } else {
                $(this).hide()
            }
        })
    }

    function hideToolWindow(object) {
        object.removeClass("active")
        let target = object.attr("target")
        $("#" + target).parent().hide()
    }

    function createToolWindowTabElement() {
        tabs.append(
            $("<div class='toolWindowButton' target='" + id + "'>" + name + "</div>")
                .click(function () {
                    if ($(this).hasClass("active")) {
                        hideToolWindow($(this))
                    } else {
                        showToolWindow(name, cssClass, position, linkedEditor, fillFunction)
                        if (linkedEditor) {
                            showEditor(linkedEditor, "")
                        }
                    }
                    setSidebarState();
                })
        )
    }

    async function createToolWindowContent() {
        let selector = $("#toolWindows ." + position)
        let content = await fillFunction
        selector.append("<div class='toolWindow " + cssClass + "' id='" + id + "'>" + content + "</div>")
    }
}

function getObjectID(s) {
    return s.toLowerCase().replaceAll("-", " ")
        .replaceAll(".", " ")
        .replaceAll("/", " ")
        .replaceAll(" ", "");
}

function setSidebarState() {
    let ActiveToolWindows = 0
    $(".toolWindowButton").each(function () {
        if ($(this).hasClass("active")) {
            ActiveToolWindows++
        }
    })
    if (ActiveToolWindows === 0) {
        $("#sidebar").hide()
        $("#file-analyzer .resizer").hide()
    } else if (ActiveToolWindows === 1) {
        $("#sidebar").show()
        $("#file-analyzer>.container>.resizer").show()
        $("#sidebar .resizer").hide()
    } else {
        $("#sidebar").show()
        $("#file-analyzer .resizer").show()
    }
}

$(document).ready(function () {
    $("#select-dir").on('click', async () => {
        var openedDir = await window.go.main.App.OpenFolder()
        if (openedDir.length > 0) {
            $("#file-uploader").hide();
            $("#file-analyzer").show();
            render()
        }
    })
    $("#select-archive").on('click', async () => {
        var openedArchive = await window.go.main.App.OpenArchive()
        if (openedArchive.length > 0) {
            $("#file-uploader").hide();
            $("#file-analyzer").show();
            render()
        }
    })
    $("#toolWindows").on('click', '#summary input:checkbox', async function () {
        await checkChildElements(this)
        var filters = {};
        $("#summary input:checkbox").each(function () {
            filters[$(this).val()] = $(this).prop('checked');
        })
        await window.go.main.App.SetFilters(filters).then(redrawEditors())
        //Group check/uncheck functionality
        async function checkChildElements(elem) {
            var checked = $(elem).prop('checked');
            var isParent = !!$(elem).closest('li').find(' > ul').length;
            if (isParent) {
                // if a parent level checkbox is changed, locate all children
                var children = $(elem).closest('li').find('ul input[type=checkbox]');
                children.prop({
                    checked
                }); // all children will have what parent has
            } else {
                console.log("is not parent")
                // if a child checkbox is changed, locate parent and all children
                var parent = $(elem).closest('ul').closest('li').find('>label input[type=checkbox]');
                var children = $(elem).closest('ul').find('input[type=checkbox]');
                if (children.filter(':checked').length === 0) {
                    // if all children are unchecked
                    parent.prop({ checked: false, indeterminate: false });
                } else if (children.length === children.filter(':checked').length) {
                    // if all children are checked
                    parent.prop({ checked: true, indeterminate: false });
                } else {
                    // if some of the children are checked
                    parent.prop({ checked: true, indeterminate: true });
                }
            }
        }
    });

    //reveal/collapse filter items on folding-icon click
    $("#toolWindows").on('click', '.group-label>.folding-icon', function () {
        let childList = $(this).parent().find("ul")
        if (childList.is(":hidden")) {
            childList.show()
            $(this).html("<svg width=\"11\" height=\"14\" viewBox=\"0 0 11 14\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n" +
                "<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M6.43427 5.93433C6.74669 5.62191 7.25322 5.62191 7.56564 5.93433C7.87806 6.24675 7.87806 6.75328 7.56564 7.0657L3.99995 10.6314L0.434266 7.0657C0.121846 6.75328 0.121846 6.24675 0.434266 5.93433C0.746685 5.62191 1.25322 5.62191 1.56564 5.93433L3.99995 8.36864L6.43427 5.93433Z\" fill=\"#6E6E6E\"/>\n" +
                "</svg>\n");
        } else {
            childList.hide()
            $(this).html("<svg width=\"11\" height=\"14\" viewBox=\"0 0 11 14\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n" +
                "<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M2.93427 4.56567C2.62185 4.25325 2.62185 3.74672 2.93427 3.4343C3.24668 3.12188 3.75322 3.12188 4.06564 3.4343L7.63132 6.99999L4.06564 10.5657C3.75322 10.8781 3.24669 10.8781 2.93427 10.5657C2.62185 10.2533 2.62185 9.74672 2.93427 9.4343L5.36858 6.99999L2.93427 4.56567Z\" fill=\"#6E6E6E\"/>\n" +
                "</svg>\n");
        }


    })
    $("#toolWindows").on('click', '.other-files li', function () {
        $(".other-files li").removeClass("active")
        $(this).addClass("active")
        let fileUUID = $(this).attr("target");
        let fileName = $(this).innerText
        let editorName = getObjectID(fileUUID)
        showEditor(editorName, window.go.main.App.GetOtherFileContent(fileUUID)).then(function () {
            let editor = ace.edit(editorName)
            editor.renderer.scrollToLine(0)
            editor.clearSelection();
        })
    });
})