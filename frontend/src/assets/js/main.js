const render = async () => {
    await clearToolWindows();
    await redrawEditors()
}
const redrawEditors = async () => {
    $("#editors>div").remove();
    await renderMainScreen();
}
const clearToolWindows = async () => {
    let tabs = $("#toolWindows-buttons").find(".toolWindowButton")
    let toolWindows = $("#toolWindows").find(".toolWindow")
    toolWindows.remove();
    tabs.remove();
}
async function renderMainScreen() {
    await showToolWindow("Filters", "filters", "top", "Main Editor", window.go.main.App.GetFilters())
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
    if (!$(`#${id}`).length) {
        if (content) await cerateEditor()
    }
    $("#editors>div").hide()
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
            showLineNumbers: false,
            showGutter: true,
            showPrintMargin: false,
            highlightSelectedWord: true,
        })
        editor.setValue(await content);
        editor.renderer.scrollToLine(Number.POSITIVE_INFINITY)
        editor.clearSelection();
        editor.session.foldAll(0, editor.session.getLength() - 4, 1);
        editor.execCommand('find');
        editor.on("click", ThreadDumpLinkHandler)
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
            } else {
                $(this).hide()
            }
        })
    }

    function createToolWindowTabElement() {
        tabs.append(
            $("<div class='toolWindowButton' target='" + id + "'>" + name + "</div>")
                .click(function () {
                    showToolWindow(name, cssClass, position, linkedEditor, fillFunction)
                    if (linkedEditor) {
                        showEditor(linkedEditor, "")
                    }
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
    $("#toolWindows").on('click', '#filters input:checkbox', async () => {
        var filters = {};
        $("#filters input:checkbox").each(function () {
            filters[$(this).val()] = $(this).prop('checked');
        })
        await window.go.main.App.SetFilters(filters)
        await redrawEditors()
    });
})

