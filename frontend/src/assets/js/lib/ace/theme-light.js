define("ace/theme/light",["require","exports","module","ace/lib/dom"], function(require, exports, module) {

    exports.isDark = false;
    exports.cssClass = "ace-github";
    exports.cssText = "";
    $.ajax({
        async: false,
        url: "/assets/js/lib/ace/theme-light.min.css",
        dataType: "text",
        success: function (cssText) {
            exports.cssText = cssText;
        }
    });

    var dom = require("../lib/dom");
    dom.importCssString(exports.cssText, exports.cssClass, false);
});                (function() {
    window.require(["ace/theme/light"], function(m) {
        if (typeof module == "object" && typeof exports == "object" && module) {
            module.exports = m;
        }
    });
})();
            