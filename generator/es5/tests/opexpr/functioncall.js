$module('functioncall', function () {
  var $static = this;
  $static.TEST = function () {
    return $g.functioncall.AnotherFunction($t.fastbox(false, $g.________testlib.basictypes.Boolean));
  };
  $static.AnotherFunction = function (param) {
    return $t.fastbox(!param.$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
