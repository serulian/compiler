$module('nativetoany', function () {
  var $static = this;
  $static.TEST = function () {
    var b;
    b = true;
    return $t.cast($t.cast(b, $t.any, true), $g.____testlib.basictypes.Boolean, false);
  };
});
