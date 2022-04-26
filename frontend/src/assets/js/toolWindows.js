const clearToolWindows = async () => {
    let tabs = $("#toolWindows-buttons").find(".toolWindowButton")
    let toolWindows = $("#toolWindows").find(".toolWindow")
    toolWindows.remove();
    tabs.remove();
}
const toolWindows = $("#toolWindows")

$(document).ready(function () {
    // Event handler for filter checkboxes
    toolWindows.on('click', '#summary input:checkbox', async function () {
        await checkChildElements(this)
        var filters = {};
        $("#summary input:checkbox").each(function () {
            filters[$(this).val()] = $(this).prop('checked');
        })
        await window.go.main.App.SetFilters(filters).then(redrawEditors())
        //Group check/uncheck functionality
        async function checkChildElements(elem) {
            var checked = $(elem).prop('checked');
            var isParent = !!$(elem).closest('li').find(' > ul').length;
            if (isParent) {
                // if a parent level checkbox is changed, locate all children
                var children = $(elem).closest('li').find('ul input[type=checkbox]');
                children.prop({
                    checked
                }); // all children will have what parent has
            } else {
                console.log("is not parent")
                // if a child checkbox is changed, locate parent and all children
                var parent = $(elem).closest('ul').closest('li').find('>label input[type=checkbox]');
                var children = $(elem).closest('ul').find('input[type=checkbox]');
                if (children.filter(':checked').length === 0) {
                    // if all children are unchecked
                    parent.prop({ checked: false, indeterminate: false });
                } else if (children.length === children.filter(':checked').length) {
                    // if all children are checked
                    parent.prop({ checked: true, indeterminate: false });
                } else {
                    // if some of the children are checked
                    parent.prop({ checked: true, indeterminate: true });
                }
            }
        }
    });

    //reveal/collapse filter items on folding-icon click
    toolWindows.on('click', '.group-label>.folding-icon', function () {

        let childList = $(this).parent().find("ul")
        if (childList.is(":hidden")) {
            childList.show()
            drawDownArrow(this)
        } else {
            childList.hide()
            drawRightArrow(this)
        }
        function drawDownArrow(el){
            $(el).html("<svg width=\"11\" height=\"14\" viewBox=\"0 0 11 14\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n" +
                "<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M6.43427 5.93433C6.74669 5.62191 7.25322 5.62191 7.56564 5.93433C7.87806 6.24675 7.87806 6.75328 7.56564 7.0657L3.99995 10.6314L0.434266 7.0657C0.121846 6.75328 0.121846 6.24675 0.434266 5.93433C0.746685 5.62191 1.25322 5.62191 1.56564 5.93433L3.99995 8.36864L6.43427 5.93433Z\" fill=\"#6E6E6E\"/>\n" +
                "</svg>\n");
        }

        function drawRightArrow(el) {
            $(el).html("<svg width=\"11\" height=\"14\" viewBox=\"0 0 11 14\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n" +
                "<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M2.93427 4.56567C2.62185 4.25325 2.62185 3.74672 2.93427 3.4343C3.24668 3.12188 3.75322 3.12188 4.06564 3.4343L7.63132 6.99999L4.06564 10.5657C3.75322 10.8781 3.24669 10.8781 2.93427 10.5657C2.62185 10.2533 2.62185 9.74672 2.93427 9.4343L5.36858 6.99999L2.93427 4.56567Z\" fill=\"#6E6E6E\"/>\n" +
                "</svg>\n");
        }

    })

    //show/hide other files on click
    toolWindows.on('click', '.other-files li', function () {
        let fileUUID = $(this).attr("target");
        let editorName = getObjectID(fileUUID)
        if (this.classList.contains("active")) {
            showEditor("Main Editor")
            $(this).removeClass("active")
        } else {
            $(".other-files li").removeClass("active")
            $(this).addClass("active")
            showEditor(editorName, window.go.main.App.GetOtherFileContent(fileUUID)).then(function () {
                let editor = ace.edit(editorName)
                editor.renderer.scrollToLine(0)
                editor.clearSelection();
            })
        }

    });
})

// showToolWindow adds tool window on the left.
// @name is the Title of Tool Window
// @position is one of top/bot
// @fillFunction is a promise for Go function that returns a string of Tool Windiow content
async function showToolWindow(name, cssClass, position, linkedEditor, fillFunction) {
    let tabs = $("#toolWindows-buttons ." + position);
    let tabsElements = tabs.find(`.toolWindowButton`);
    let toolWindows = $("#toolWindows ." + position + " .toolWindow")
    let id = getObjectID(name);
    if (!getToolWindowTabElement()) {
        createToolWindowTabElement()
        await createToolWindowContent()
        return showToolWindow(name, cssClass, position, linkedEditor, fillFunction)
    }
    selectToolWindowTab()

    function getToolWindowTabElement() {
        let a;
        tabsElements.each(function () {
            let attr = $(this).attr('target')
            if (attr === id) {
                a = $(this)
                return a
            }
        })
        return a
    }

    function selectToolWindowTab() {
        tabsElements.each(function () {
            if ($(this).attr('target') === id) {
                $(this).addClass("active")
            } else {
                $(this).removeClass("active")
            }
        })
        showToolWindowContent()
    }

    function showToolWindowContent() {
        toolWindows.each(function () {
            if ($(this).prop('id') === id) {
                $(this).show()
                $(this).parent().show()
            } else {
                $(this).hide()
            }
        })
    }

    function hideToolWindow(object) {
        object.removeClass("active")
        let target = object.attr("target")
        $("#" + target).parent().hide()
    }

    function removeToolWindow(object) {
        let target = object.attr("target")
        $("#" + target).remove()
        object.remove()
    }

    function createToolWindowTabElement() {
        let closeButton = function () {
            if (id !== getObjectID("Summary") && id !== getObjectID("Static Info")) {
                return '<span class="closebtn">&times;</span>'
            }
            return ''
        }
        tabs.append(
            $("<div class='toolWindowButton' target='" + id + "'>" + name + closeButton() + '</div>')
                .click(function (e) {
                    if ($(e.target).hasClass("closebtn")) {
                        removeToolWindow($(this))
                        renderMainScreen()
                    } else if ($(this).hasClass("active")) {
                        hideToolWindow($(this))
                    } else {
                        showToolWindow(name, cssClass, position, linkedEditor, fillFunction)
                        if (linkedEditor) {
                            showEditor(linkedEditor, "")
                        }
                    }
                    setSidebarState();
                })
        )
    }

    async function createToolWindowContent() {
        let selector = $("#toolWindows ." + position)
        let content = await fillFunction
        selector.append("<div class='toolWindow " + cssClass + "' id='" + id + "'>" + content + "</div>")
    }
}

function setSidebarState() {
    let ActiveToolWindows = 0
    $(".toolWindowButton").each(function () {
        if ($(this).hasClass("active")) {
            ActiveToolWindows++
        }
    })
    if (ActiveToolWindows === 0) {
        $("#sidebar").hide()
        $("#file-analyzer .resizer").hide()
    } else if (ActiveToolWindows === 1) {
        $("#sidebar").show()
        $("#file-analyzer>.container>.resizer").show()
        $("#sidebar .resizer").hide()
    } else {
        $("#sidebar").show()
        $("#file-analyzer .resizer").show()
    }
}