function setPreferableColorScheme() {
    getPreferredColorScheme().then(scheme => {
        setColorScheme(scheme);
    });
}
function setColorScheme(scheme) {
    switch (scheme) {
        case 'dark':
            document.body.classList.add("dark");
            break
        case 'light':
            document.body.classList.remove("dark");
            break
    }
}

async function getPreferredColorScheme() {
    let theme = await window.go.main.App.GetSetting("EditorTheme")
    if (theme!=="system") {
        return theme;
    } else {
        return getSystemTheme();
    }
    function getSystemTheme() {
        if (window.matchMedia) {
            if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                return 'dark';
            } else {
                return 'light';
            }
        }
        return 'light';
    }
}

if (window.matchMedia) {
    //set theme on startup
    setPreferableColorScheme();
    //add event listener for theme change
    var colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');
    colorSchemeQuery.addEventListener('change', function (){
        setPreferableColorScheme();
    });
}