define("ace/theme/idealog",["require","exports","module","ace/lib/dom"], function(require, exports, module) {

    exports.isDark = false;
    exports.cssClass = "editor";
    exports.cssText = "";
    $.ajax({
        async: false,
        url: "/assets/css/editor/theme-idealog.min.css",
        dataType: "text",
        success: function (cssText) {
            exports.cssText = cssText;
        }
    });

    var dom = require("../lib/dom");
    dom.importCssString(exports.cssText, exports.cssClass, false);
});                (function() {
    window.require(["ace/theme/idealog"], function(m) {
        if (typeof module == "object" && typeof exports == "object" && module) {
            module.exports = m;
        }
    });
})();
            