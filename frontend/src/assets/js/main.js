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

//Show/hide options list on dropdown click
    $(document).on("click", ".dropdown", function () {
        $(this).find('.options').first().toggle()
    })
})
