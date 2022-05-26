document.addEventListener('DOMContentLoaded', function () {
    let settingsOverlay;

    function applySettingsForAllEditors(settings) {
        console.log(settings);
        editors.find(".editor").each(function () {
            let editor = ace.edit(this.attributes["id"].value);
            editor.setFontSize(settings.EditorFontSize);
            editor.session.setUseWrapMode(settings.EditorDefaultSoftWrapState);
        });
    }

    window.runtime.EventsOn("ShowSettings", async function () {
        await ShowSettings()
        settingsOverlay = $('#settings-overlay');
        settingsOverlay.on('click', function (e) {
            if (e.target.id === settingsOverlay.attr('id') || e.target.className === "settings-overlay-close") {
                settingsOverlay.remove();
            }
        });
    });
    window.runtime.EventsOn("SettingsChanged", function (settings) {
        window.runtime.LogDebug("Settings changed: " + JSON.stringify(settings));
        applySettingsForAllEditors(settings);
        setPreferableColorScheme();
    });

    $(document).keydown(function (e) {
        if ((e.key === "Escape") && settingsOverlay.is(":visible")) {
            e.preventDefault();
            settingsOverlay.remove();
        }
    });
});

async function ShowSettings() {
    await window.go.main.App.GetSettingsScreenHTML().then(function (data) {
        $("body").append(data);
    });
}

function SaveSetting(id, option) {
    window.go.main.App.SaveSetting(id, option);
}

$(document).on("DOMSubtreeModified", ".settings-section-content-item .dropdown", function () {
    let settingId = $(this).attr("setting")
    let settingValue = $(this).find("li.active").first().attr("value")
    if (parseInt(settingValue)) {
        settingValue = parseInt(settingValue)
    }
    window.go.main.App.SaveSetting(settingId, settingValue);
})