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
                regex: /^$/,
                token: "empty_line"
            }, {
                regex: /\s—\s(.*?)\s—\s/,
                token: "variable.class"
            }, {
                regex: /INFO|INDEX|SEVERE|VERB|TRACE/,
                token: "loglevel.info",
            },{
                regex: /ERROR|PARSE_ERROR|FREEZE|STDERR/,
                token: "loglevel.error",
            },{
                regex: /WARN|STDERR/,
                token: "loglevel.warn",
            }, {
                regex: /(threadDump\S*(?=\s)*)/,
                token: "ThreadDumpsHyperlink",
            },{
                regex: /(Indexing project:.*)(report.html)(.*Report: )(.*\.html)/,
                token: ["text", "IndexingProjectDiagnosticHyperlink", "text", "IndexingDiagnosticHyperlink"],
            }, {
                regex: /(\d{2}\s(Jan|JAN|Feb|FEB|Mar|MAR|Apr|APR|May|MAY|Jun|JUN|Jul|JUL|Aug|AUG|Sep|SEP|Oct|OCT|Nov|NOV|Dec|DEC)\s\d{4}\s\d{1,2}:\d{2}:\d{2})|\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}[,|:]\d{3}/,
                token: "date"
            }, {
                regex: /^\s*at.*$|STDERR.*at\s.*$/,
                token: "loglevel.warn"
            }, {
                defaultToken: "text"
            }]
        };
        this.normalizeRules()
    };

    oop.inherits(CustomHighlightRules, TextHighlightRules);

    exports.CustomHighlightRules = CustomHighlightRules;
});

//Folds IDE STARTED ... IDE SHUTDOWN sections
define("ace/mode/folding/idea-style", [], function (require, exports, module) {
    "use strict";

    let oop = require("../../lib/oop");
    let Range = require("../../range").Range;
    let BaseFoldMode = require("./fold_mode").FoldMode;

    let FoldMode = exports.FoldMode = function (commentRegex) {
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

    (function () {
        this.foldingStartMarker = /(\{|\[)[^\}\]]*$|^\s*(\/\*)/;
        this.foldingStopMarker = /^[^\[\{]*(\}|\])|^[\s\*]*(\*\/)/;
        this.singleLineBlockCommentRe = /^\s*(\/\*).*\*\/\s*$/;
        this.tripleStarBlockCommentRe = /^\s*(\/\*\*\*).*\*\/\s*$/;
        this.startRegionRe = /-+ IDE STARTED -+/;
        this._getFoldWidgetBase = this.getFoldWidget;
        this.getFoldWidget = function (session, foldStyle, row) {
            var line = session.getLine(row);

            if (this.singleLineBlockCommentRe.test(line)) {
                if (!this.startRegionRe.test(line) && !this.tripleStarBlockCommentRe.test(line))
                    return "";
            }

            var fw = this._getFoldWidgetBase(session, foldStyle, row);

            if (!fw && this.startRegionRe.test(line) || row === 0)
                return "start"; // lineCommentRegionStart

            return fw;
        };
        this.getFoldWidgetRange = function (session, foldStyle, row, forceMultiline) {
            var line = session.getLine(row);
            if (this.startRegionRe.test(line) || row === 0)
                return this.getCommentRegionBlock(session, line, row);
        };
        this.getCommentRegionBlock = function (session, line, row) {
            let startColumn = line.search(this.startRegionRe);
            let maxRow = session.getLength();
            let startRow = row;

            let lineIdeShutdown = /-+ IDE SHUTDOWN -+/;
            let lineWebserverStopped = /.*web server stopped.*/;
            let depth = 1;
            while (++row < maxRow) {
                line = session.getLine(row);
                let lineMatchIdeShutdown = lineIdeShutdown.exec(line);
                let lineMatchIdeStart = this.startRegionRe.test(line);
                if (lineMatchIdeStart) {
                    depth--;
                    row = row - 1
                    line = session.getLine(row);
                }
                if (lineMatchIdeShutdown && lineMatchIdeShutdown[0]) {
                    depth--;
                    if (lineWebserverStopped.exec(session.getLine(row + 1))) {
                        row = row + 1
                        line = session.getLine(row);
                    }
                }
                if (!depth) break;
            }

            let endRow = row;
            if (endRow > startRow && depth===0) {
                return new Range(startRow, startColumn, endRow, line.length);
            }
        };
    }).call(FoldMode.prototype);
});