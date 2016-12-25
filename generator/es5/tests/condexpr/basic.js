$module('basic', function () {
  var $static = this;
  $static.TEST = function () {
    var value;
    value = $t.fastbox(2, $g.____testlib.basictypes.Integer);
    return value.$wrapped == 2 ? $t.fastbox(true, $g.____testlib.basictypes.Boolean) : $t.fastbox(false, $g.____testlib.basictypes.Boolean);
  };
});
