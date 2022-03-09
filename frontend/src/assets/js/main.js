const filters = $("#filters");
const selectarchive = document.getElementById('select-archive');
const selectdir = document.getElementById('select-dir');

const render = async () => {
    var editor = ace.edit("editor");
    editor.setOptions({
        mode: 'ace/mode/idea_log',
        theme: "ace/theme/light",
        readOnly: true,
        selectionStyle: "text",
        showLineNumbers: false,
        showGutter: false,
        showPrintMargin: false,
        highlightSelectedWord: true,
    })
    filters[0].innerHTML = await window.go.main.App.GetFilters()
    editor.setValue(await window.go.main.App.GetLogs());
    editor.renderer.scrollToLine(Number.POSITIVE_INFINITY)
    editor.clearSelection();
    editor.execCommand('find');
};


$(document).ready(function () {
    selectdir.onclick = async () => {
        var openedDir = await window.go.main.App.OpenFolder()
        if (openedDir.length > 0) {
            $("#file-uploader").hide();
            $("#file-analyzer").show();
            render()
        }
    }
    selectarchive.onclick = async () => {
        var openedArchive = await window.go.main.App.OpenArchive()
        if (openedArchive.length > 0) {
            $("#file-uploader").hide();
            $("#file-analyzer").show();
            render()
        }
    }
    filters.on('click', 'input:checkbox', async () => {
        var filters = {};
        $("#filters input:checkbox").each(function () {
            filters[$(this).val()] = $(this).prop('checked');
        })
        await window.go.main.App.SetFilters(filters)
        await render()
    });

})

