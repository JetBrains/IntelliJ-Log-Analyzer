function getObjectID(s) {
    return s.toLowerCase().replaceAll("-", " ")
        .replaceAll(".", " ")
        .replaceAll("/", " ")
        .replaceAll(" ", "");
}

// Dropdown element logic
$(document).ready(function () {
    function setDropDownTitle(e) {
        let dropdown = $(e).closest(".dropdown")
        dropdown.find('.title').first().html(dropdown.find("li.active").first().html())
    }

    // Set dropdown title from active component in options
    $(document).on("DOMSubtreeModified", ".dropdown .options", function () {
        setDropDownTitle(this)
    })

    //Dropdown set active li element
    $(document).on("click", ".dropdown li", function () {
        $(this).siblings().each(function () {
            $(this).removeClass("active")
        })
        $(this).addClass("active")
        setDropDownTitle(this)
    })

    //Show/hide options list on click. Also closes all but current dropdowns
    $(document).on("click", function (e) {
        $(this).find('.dropdown .options:visible').each(function () {
            if (!$(this).siblings(".dropdown .title").is($(e.target))) {
                $(this).hide()
            }
        })
        if (e.target.classList.contains("title") || e.target.classList.contains("dropdown")) {
            $(e.target).closest(".dropdown").find('.options').first().toggle()
        }
    })
})
//Program-wide key bindings
$(document).keydown(function (e) {
    if ((e.key === "f") && (e.ctrlKey||e.metaKey)) {
        e.preventDefault();
        $("#editors>div").find(".search-box:visible input").focus();
    }
});