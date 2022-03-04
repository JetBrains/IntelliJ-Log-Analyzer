define('ace/mode/idea_log', [], function(require, exports, module) {
    var oop = require("ace/lib/oop");
    var TextMode = require("ace/mode/text").Mode;
    var Tokenizer = require("ace/tokenizer").Tokenizer;
    var CustomHighlightRules = require("ace/mode/idea_log_highlight_rules").CustomHighlightRules;

    var Mode = function() {
        this.HighlightRules = CustomHighlightRules;
    };
    oop.inherits(Mode, TextMode);

    (function() {

    }).call(Mode.prototype);

    exports.Mode = Mode;
});

define('ace/mode/idea_log_highlight_rules', [], function(require, exports, module) {
    var oop = require("ace/lib/oop");
    var TextHighlightRules = require("ace/mode/text_highlight_rules").TextHighlightRules;

    var CustomHighlightRules = function() {

        var keywordMapper = this.createKeywordMapper({
            "variable.language": "this",
            "keyword": "Mark|Ben|Bill",
            "constant.language": "true|false|null",
            // it is also possible to use css, but that may conflict with themes
            // "customTokenName": "problem"
        }, "text", true);

        this.$rules = {
            "start" : [{
                token : "empty_line",
                regex : '^$'
            },{
                regex: "INFO",
                token: "support.constant",
            },{
                regex: "WARN",
                token: "support.constant",
            },{
                regex: "\\d{2} (Jan|JAN|Feb|FEB|Mar|MAR|Apr|APR|May|MAY|Jun|JUN|Jul|JUL|Aug|AUG|Sep|SEP|Oct|OCT|Dec|DEC) \\d{4} \\d{1,2}:\\d{2}:\\d{2}",
                token: "keyword.control"
            },{
                defaultToken : "text"
            }]
        };
        this.normalizeRules()
    };

    oop.inherits(CustomHighlightRules, TextHighlightRules);

    exports.CustomHighlightRules = CustomHighlightRules;
});