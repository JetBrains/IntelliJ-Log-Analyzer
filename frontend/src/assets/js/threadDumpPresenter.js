define("ace/ext/threadDumpPresenter", [],function() {
    let handler = async function (e) {
        let editor = e.editor
        let pos = editor.getCursorPosition()
        let token = editor.session.getTokenAt(pos.row, pos.column)
        console.log(token)
        if ((token.type !== null) && (/hyperlink/.test(token.type))) {
            await openThreadDump(token.value)
        }

    }
    ace.edit("main-editor").on("click", handler)
});

async function openThreadDump(path) {
    showToolWindow("ThreadDump","top","ThreadDump Editor",window.go.main.App.GetThreadDumpsFilters())
    showEditor("ThreadDump Editor", window.go.main.App.GetThreadDumps(path));
}