define('ace/theme/hangulize', [
  'require',
  'exports',
  'module',
  'ace/lib/dom',
], function(require, exports, module) {

exports.isDark = false;
exports.cssClass = 'ace-hangulize';
exports.cssText = '\
\
.ace-hangulize {\
  background-color: #fff;\
  color: #f00\
}\
.ace-hangulize .ace_cursor {\
  color: #432\
}\
.ace-hangulize .ace_marker-layer .ace_selection {\
  background: #def\
}\
.ace-hangulize .ace_marker-layer .ace_active-line {\
  background: #ffd\
}\
.ace-hangulize .ace_marker-layer .ace_selected-word {\
  border: 2px solid #9af\
}\
.ace-hangulize .ace_keyword {\
  color: #89a\
}\
.ace-hangulize .ace_invalid {\
  color: #fff;\
  background-color: #fa9\
}\
.ace-hangulize .ace_storage {\
  color: #a94\
}\
.ace-hangulize .ace_variable {\
  color: #39c\
}\
.ace-hangulize .ace_string {\
  color: #432\
}\
.ace-hangulize .ace_string.ace_double {\
  color: #39c\
}\
.ace-hangulize .ace_string.ace_interpolated {\
  color: #b8c\
}\
.ace-hangulize .ace_hangul.ace_choseong {\
  color: #c24;\
}\
.ace-hangulize .ace_hangul.ace_jungseong {\
  color: #791\
}\
.ace-hangulize .ace_hangul.ace_jongseong {\
  color: #36c;\
}\
.ace-hangulize .ace_comment {\
  color: #44f;\
}\
\
';

let dom = require('ace/lib/dom');
dom.importCssString(exports.cssText, exports.cssClass);

});

(function() {
  window.require(['ace/theme/hangulize'], function(m) {
    if (typeof module == 'object' && typeof exports == 'object' && module) {
      module.exports = m;
    }
  });
})();