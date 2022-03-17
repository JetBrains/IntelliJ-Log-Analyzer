let ThreadDumpLinkHandler = async function (e) {
    let editor = e.editor
    let pos = editor.getCursorPosition()
    let token = editor.session.getTokenAt(pos.row, pos.column)
    console.log(token)
    if ((token.type !== null) && (/hyperlink/.test(token.type))) {
        await openThreadDump(token.value)
    }

}

let iterator = 1;
async function openThreadDump(path) {
    var myRegexp = new RegExp("(\\d{8}-)(\\d{6})", "g");
    var match = myRegexp.exec(path);
    let name = "TD-" + match[2]
    let editorName = "threadDump editor" + path.toLowerCase()
                                            .replaceAll(" ","")
                                            .replaceAll("-"," ")
                                            .replaceAll("."," ")
                                            .replaceAll("/"," ");
    showToolWindow(name, "top", editorName, window.go.main.App.GetThreadDumpsFilters())
    await showEditor(editorName, window.go.main.App.GetThreadDumps(path));
    iterator++
}