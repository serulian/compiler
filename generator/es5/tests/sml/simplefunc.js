$module('simplefunc', function () {
  var $static = this;
  $static.SimpleFunction = function () {
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.simplefunc.SimpleFunction();
  };
});
