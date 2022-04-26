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
    if (await window.go.main.App.GetStaticInfo()) {
        await showToolWindow("Static Info", "staticinfo", "bot", "", window.go.main.App.GetStaticInfo())
    }
    window.mainEditorID = showEditor("Main Editor", window.go.main.App.GetLogs());
    setSidebarState();
    function setSummaryToolWindowGroupCheckboxStates() {
        $("#summary input:checkbox").each(function (){
            if ($(this).attr("mixed")!==undefined) {
                $(this).prop({indeterminate: true});
            }
        })
    }
}