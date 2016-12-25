$module('basic', function () {
  var $static = this;
  $static.DoSomething = function () {
    return $t.fastbox(42, $g.____testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    return $t.fastbox($g.basic.DoSomething().$wrapped == 42, $g.____testlib.basictypes.Boolean);
  };
});
