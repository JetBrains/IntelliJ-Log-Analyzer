define('ace/mode/idea_log', [], function (require, exports, module) {
    let oop = require("ace/lib/oop");
    let TextMode = require("ace/mode/text").Mode;
    let IdeaStyleFoldMode = require("./folding/idea-style").FoldMode;
    require("ace/ext/threadDumpPresenter");
    let Tokenizer = require("ace/tokenizer").Tokenizer;
    let CustomHighlightRules = require("ace/mode/idea_log_highlight_rules").CustomHighlightRules;

    let Mode = function () {
        this.HighlightRules = CustomHighlightRules;
        this.foldingRules = new IdeaStyleFoldMode();
    };
    oop.inherits(Mode, TextMode);

    (function () {

    }).call(Mode.prototype);

    exports.Mode = Mode;
});

define('ace/mode/idea_log_highlight_rules', [], function (require, exports, module) {
    var oop = require("ace/lib/oop");
    var TextHighlightRules = require("ace/mode/text_highlight_rules").TextHighlightRules;

    var CustomHighlightRules = function () {

        this.$rules = {
            "start": [{
                token: "empty_line",
                regex: '^$'
            }, {
                regex: "INFO",
                token: "loglevel.info",
            }, {
                regex: "WARN",
                token: "loglevel.warn",
            }, {
                regex: "ERROR",
                token: "loglevel.error",
            },{
                regex: "FREEZE",
                token: "loglevel.error",
            },{
                regex: "(threadDump\\S*(?=\\s)*)",
                token: "hyperlink",
            }, {
                regex: "\\d{2} (Jan|JAN|Feb|FEB|Mar|MAR|Apr|APR|May|MAY|Jun|JUN|Jul|JUL|Aug|AUG|Sep|SEP|Oct|OCT|Dec|DEC) \\d{4} \\d{1,2}:\\d{2}:\\d{2}",
                token: "date"
            }, {
                regex: "^\\s*at.*$",
                token: "loglevel.warn"
            }, {
                regex: " - (.*) - ",
                token: "variable.class"
            }, {
                    defaultToken: "text"
            }]
        };
        this.normalizeRules()
    };

    oop.inherits(CustomHighlightRules, TextHighlightRules);

    exports.CustomHighlightRules = CustomHighlightRules;
});

define("ace/mode/folding/idea-style",[], function(require, exports, module) {
    "use strict";

    let oop = require("../../lib/oop");
    let Range = require("../../range").Range;
    let BaseFoldMode = require("./fold_mode").FoldMode;

    let FoldMode = exports.FoldMode = function(commentRegex) {
        if (commentRegex) {
            this.foldingStartMarker = new RegExp(
                this.foldingStartMarker.source.replace(/\|[^|]*?$/, "|" + commentRegex.start)
            );
            this.foldingStopMarker = new RegExp(
                this.foldingStopMarker.source.replace(/\|[^|]*?$/, "|" + commentRegex.end)
            );
        }
    };
    oop.inherits(FoldMode, BaseFoldMode);

    (function() {
        this.foldingStartMarker = /(\{|\[)[^\}\]]*$|^\s*(\/\*)/;
        this.foldingStopMarker = /^[^\[\{]*(\}|\])|^[\s\*]*(\*\/)/;
        this.singleLineBlockCommentRe= /^\s*(\/\*).*\*\/\s*$/;
        this.tripleStarBlockCommentRe = /^\s*(\/\*\*\*).*\*\/\s*$/;
        this.startRegionRe = /-+ IDE STARTED -+/;
        this._getFoldWidgetBase = this.getFoldWidget;
        this.getFoldWidget = function(session, foldStyle, row) {
            var line = session.getLine(row);

            if (this.singleLineBlockCommentRe.test(line)) {
                if (!this.startRegionRe.test(line) && !this.tripleStarBlockCommentRe.test(line))
                    return "";
            }

            var fw = this._getFoldWidgetBase(session, foldStyle, row);

            if (!fw && this.startRegionRe.test(line))
                return "start"; // lineCommentRegionStart

            return fw;
        };
        this.getFoldWidgetRange = function(session, foldStyle, row, forceMultiline) {
            var line = session.getLine(row);
            if (this.startRegionRe.test(line))
                return this.getCommentRegionBlock(session, line, row);
        };
        this.getCommentRegionBlock = function(session, line, row) {
            let startColumn = line.search(this.startRegionRe);
            let maxRow = session.getLength();
            let startRow = row;

            let re = /.*--- IDE SHUTDOWN ---.*/;
            let reWebserverStopped = /.*web server stopped.*/;
            let depth = 1;
            while (++row < maxRow) {
                line = session.getLine(row);
                let m = re.exec(line);

                if (!m) continue;
                 if (m[0]) {
                     depth--;
                     if (reWebserverStopped.exec(session.getLine(row+1))) {
                         row = row+1
                     }
                 }
                else depth++;

                if (!depth) break;
            }

            let endRow = row;
            if (endRow > startRow) {
                return new Range(startRow, startColumn, endRow, line.length);
            }
        };
    }).call(FoldMode.prototype);
});