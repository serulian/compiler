$module('plus', function () {
  var $static = this;
  $static.SomeFunction = function () {
    return $t.fastbox(1 + 2, $g.________testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    return $t.fastbox($g.plus.SomeFunction().$wrapped == 3, $g.________testlib.basictypes.Boolean);
  };
});
