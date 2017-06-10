$module('plus', function () {
  var $static = this;
  $static.SomeFunction = function () {
    return $t.fastbox(1 + 2, $g.____testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    return $t.fastbox($g.plus.SomeFunction().$wrapped == 3, $g.____testlib.basictypes.Boolean);
  };
});
