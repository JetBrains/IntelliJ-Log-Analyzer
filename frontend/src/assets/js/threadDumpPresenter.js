let ThreadDumpLinkHandler = async function (e) {
    let editor = e.editor
    let pos = editor.getCursorPosition()
    let token = editor.session.getTokenAt(pos.row, pos.column)
    if ((token.type !== null) && (/hyperlink/.test(token.type))) {
        await openThreadDump(token.value)
    }

}

async function openThreadDump(path) {
    var myRegexp = new RegExp("(\\d{8}-)(\\d{6})", "g");
    var match = myRegexp.exec(path);
    let name;
    if (match && match[2]) {
        name = "TD-" + match[2]
    } else {
        name = path
    }
    let id = getObjectID(name);
    let cssClass = "ThreadDumpFilter"
    let editorName = getObjectID("threadDump editor" + path.toLowerCase());
    await showToolWindow(name, cssClass, "top", editorName, window.go.main.App.GetThreadDumpsFilters(path))
    let files = $("#" + id).children()
    files.bind('click', async function () {
        let filename = $(this).attr("filename");
        files.removeClass("active")
        $(this).addClass("active")
        await showEditor(editorName, window.go.main.App.GetThreadDumpFileContent(path, filename))
        let editor = ace.edit(editorName);
        editor.setValue(await window.go.main.App.GetThreadDumpFileContent(path, filename))
        editor.renderer.scrollToLine(0)
        editor.clearSelection();
    })
    files.first().click();
}