$module('functioncall', function () {
  var $static = this;
  $static.TEST = function () {
    return $g.functioncall.AnotherFunction($t.fastbox(false, $g.____testlib.basictypes.Boolean));
  };
  $static.AnotherFunction = function (param) {
    return $t.fastbox(!param.$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
