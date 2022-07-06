//Clear Tool Windows, and redraw editors
const render = async () => {
    await clearToolWindows();
    await redrawEditors()
    $(".group-label>.folding-icon").click();
}

//Remove all editors and get Summary screen from server
const redrawEditors = async () => {
    $("#editors>div").remove();
    await renderMainScreen();
}

//Get Summary Screen from server
async function renderMainScreen() {
    await showToolWindow("Summary", "filters", "top", "Main Editor", window.go.main.App.GetSummary())
    setSummaryToolWindowGroupCheckboxStates()
    addSummaryToolWindowListeners()
    if (await window.go.main.App.GetStaticInfo()) {
        await showToolWindow("Static Info", "staticinfo", "bot", "", window.go.main.App.GetStaticInfo())
    }
    window.mainEditorID = await showEditor("Main Editor", window.go.main.App.GetLogs());
    setSidebarState();
    function addSummaryToolWindowListeners() {
        $("#summary .link.show-in-editor").on("click",async function (e){
            e.preventDefault()

            const entityInstanceID = $(this.closest("label")).attr("for")
            const FirstInstanceString = await window.go.main.App.GetEntityInstanceFirstString(entityInstanceID)
            focusStringByContent(FirstInstanceString)

        })
    }
    function focusStringByContent(FirstInstanceString) {
        let editor = ace.edit(window.mainEditorID)
        for (let i = FirstInstanceString.length-10; i > 0; i--) {
            let stringpart = FirstInstanceString.substring(i,FirstInstanceString.length)
            let ranges = editor.findAll(stringpart,{
                wrap: true,
                caseSensitive: true,
                wholeWord: true,
                regExp: false,
            })
            if (ranges===1){
                const range =  editor.getSelection().getAllRanges();
                editor.gotoLine(range[0].end.row, 0, true);
                break;
            }
        }
        return null
    }
    function setSummaryToolWindowGroupCheckboxStates() {
        $("#summary input:checkbox").each(function (){
            if ($(this).attr("mixed")!==undefined) {
                $(this).prop({indeterminate: true});
            }
        })
    }
}