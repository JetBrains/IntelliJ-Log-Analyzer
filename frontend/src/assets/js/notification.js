function showNotification(type, text) {
    let cssClass = "info"
    if (type.length !== 0) {
        cssClass = type
    }

    $("#file-analyzer #log-holder #alerts").append(`
    <div class="alert ` + cssClass + `">
        ` + text + `
        <span class="closebtn" onclick="this.parentElement.style.display='none';">&times;</span>
    </div>
    `)
    $(".alert").delay(3000).fadeOut(0)
    //todo remove(instead of hide) alerts on close/click
}
