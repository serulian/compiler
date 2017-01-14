$module('window', function () {
  var $static = this;
  $static.TEST = function () {
    return $global.document;
  };
});
