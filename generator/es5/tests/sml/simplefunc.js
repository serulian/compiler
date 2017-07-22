$module('simplefunc', function () {
  var $static = this;
  $static.SimpleFunction = function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.simplefunc.SimpleFunction();
  };
});
