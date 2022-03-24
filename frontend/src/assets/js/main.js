const render = async () => {
    await clearToolWindows();
    await redrawEditors()
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
    await showToolWindow("Filters", "filters", "top", "Main Editor", window.go.main.App.GetSummary())
    if (await window.go.main.App.GetStaticInfo()) {
        await showToolWindow("Static Info", "staticinfo", "bot", "", window.go.main.App.GetStaticInfo())
    }
    showEditor("Main Editor", window.go.main.App.GetLogs());
}

//showEditor dhows editor if it exists, or generate new editor if it does not exist
// @name is the id attribute for the editor
// @content is a async function which returns content to be displayed
async function showEditor(name, content) {
    let id = getObjectID(name);
    let editors = $("#editors")
    $("#editors>div").hide()
    if (!$(`#${id}`).length) {
        if (content) await cerateEditor()
    }
    editors.find(`.${id}`).show()
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
        })
        editor.setValue(await content);
        await highlightEntriesTypes();
        editor.clearSelection();
        editor.execCommand('find');
        editor.renderer.scrollToLine(Number.POSITIVE_INFINITY)
        editor.on("click", ThreadDumpLinkHandler)
        editor.session.foldAll(0, editor.session.getLength() - 4, 1);

        //Checks entryType of every line and highlight this line according to type.
        //Highlighting color is configured for every DynamicEntity on init()
        async function highlightEntriesTypes() {
            let needle = /(^\s*<entryType>)(.*)(<\/entryType>\s*)/
            let mappedColors = JSON.parse(await window.go.main.App.GetEntityNamesWithLineHighlightingColors())
            editor.$search.setOptions({
                needle: needle,
                caseSensitive: true,
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
                    editor.session.addMarker(range[rangeKey], getObjectID(groupName), "screenLine", false)
                }
                editor.session.replace(range[rangeKey], "")
            }
        }
        function addCssClass(className, content) {
            document.body.appendChild(
                Object.assign(
                    document.createElement("style"),
                    {textContent: "." + className + " {" +content+"}"})
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
    return s.toLowerCase().replaceAll("-"," ")
        .replaceAll("."," ")
        .replaceAll("/"," ")
        .replaceAll(" ", "");
}
function setSidebarState() {
    let ActiveToolWindows = 0
    $(".toolWindowButton").each(function () {
        if ($(this).hasClass("active")) {
            ActiveToolWindows++
        }
    })
    if (ActiveToolWindows===0) {
        $("#sidebar").hide()
        $("#file-analyzer .resizer").hide()
    } else if (ActiveToolWindows===1) {
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
    $("#toolWindows").on('click', '#filters input:checkbox', function () {
        checkChildElements(this)
        var filters = {};
        $("#filters input:checkbox").each(function () {
            filters[$(this).val()] = $(this).prop('checked');
        })
        window.go.main.App.SetFilters(filters).then(redrawEditors())
    });

    //Group check/uncheck functionality
    function checkChildElements(elem) {
        console.log("change")
        var checked = $(elem).prop("checked"),
            container = $(elem).parent();

        container.find('input[type="checkbox"]').prop({
            indeterminate: false,
            checked: checked
        });
        console.log(container)
        checkSiblings(container);
        function checkSiblings(el) {

            var parent = el.parent().parent(),
                all = true;
            el.siblings().each(function() {
                let returnValue = all = ($(elem).children('input[type="checkbox"]').prop("checked") === checked);
                return returnValue;
            });

            if (all && checked) {

                parent.children('input[type="checkbox"]').prop({
                    indeterminate: false,
                    checked: checked
                });

                checkSiblings(parent);

            } else if (all && !checked) {

                parent.children('input[type="checkbox"]').prop("checked", checked);
                parent.children('input[type="checkbox"]').prop("indeterminate", (parent.find('input[type="checkbox"]:checked').length > 0));
                checkSiblings(parent);

            } else {

                el.parents("li").children('input[type="checkbox"]').prop({
                    indeterminate: true,
                    checked: false
                });

            }

        }

    }

})

