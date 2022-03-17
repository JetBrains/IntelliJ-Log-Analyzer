
const render = async () => {
    clearToolWindows();
    clearEditors();
    await renderMainScreen();
    function clearToolWindows() {
        let tabs = $("#toolWindows-buttons").find(".toolWindowButton")
        let toolWindows = $("#toolWindows").find(".toolWindow")
        toolWindows.remove();
        tabs.remove();
    }
    function clearEditors() {
        $("#editors>div").remove();
    }
}
function renderMainScreen() {
    showToolWindow("Filters", "top", "Main Editor", window.go.main.App.GetFilters())
    showToolWindow("Static Info", "bot", "", window.go.main.App.GetStaticInfo())
    showEditor("Main Editor", window.go.main.App.GetLogs());
}

async function showEditor(name, content) {
    let id  = name.toLowerCase().replace(" ","-");
    let editors = $("#editors")
    if (!$(`#${id}`).length) {
        if (content) await cerateEditor()
    }
    $("#editors>div").hide()
    editors.find(`.${id}`).show()
    async function cerateEditor() {
        editors.append(
            `    
            <div class=${id}>
                <div class="search-box" linked-editor="${id}"></div>
                <div id="${id}" class="editor">
                    <div class="loader">Loading...</div>
                </div>
            </div>
            `
        )
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
    }

}

// showToolWindow adds tool window on the left.
// @name is the Title of Tool Window
// @position is one of top/bot
// @fillFunction is a promise for Go function that returns a string of Tool Windiow content
function showToolWindow(name, position, linkededitor, fillFunction) {
    let tabs = $("#toolWindows-buttons ." + position);
    let tabsElements = tabs.find(`.toolWindowButton`);
    let toolWindows = $("#toolWindows ." + position + " .toolWindow")
    let id = name.toLowerCase().replace(" ", "-");
    if (!getToolWindowTabElement()) {
        createToolWindowTabElement()
        createToolWindowContent()
        showToolWindow(name, position, fillFunction)
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
                    showToolWindow(name, position, fillFunction)
                    if (linkededitor) {showEditor(linkededitor,"")}
                })
        )
    }

    async function createToolWindowContent() {
        let selector = $("#toolWindows ." + position)
        let content = await fillFunction
        selector.append("<div class='toolWindow' id='" + id + "'>" + content + "</div>")
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
    $("#toolWindows").on('click', '#filters input:checkbox', async () => {
        var filters = {};
        $("#filters input:checkbox").each(function () {
            filters[$(this).val()] = $(this).prop('checked');
        })
        await window.go.main.App.SetFilters(filters)
        await render()
    });
})

