function setColorScheme(scheme) {
    console.log("setColorScheme: " + scheme);
    switch (scheme) {
        case 'dark':
            document.body.classList.add("dark");
            break
        case 'light':
            document.body.classList.remove("dark");
            break
    }
}

function getPreferredColorScheme() {
    if (window.matchMedia) {
        if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
            return 'dark';
        } else {
            return 'light';
        }
    }
    return 'light';
}

if (window.matchMedia) {
    //set theme on startup
    setColorScheme(getPreferredColorScheme());
    //add event listener for theme change
    var colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');
    colorSchemeQuery.addEventListener('change', function (){
        setColorScheme(getPreferredColorScheme());
    });
}