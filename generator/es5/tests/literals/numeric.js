$module('numeric', function () {
  var $static = this;
  $static.DoSomething = function () {
    $t.fastbox(2, $g.____testlib.basictypes.Integer);
    $t.fastbox(2.0, $g.____testlib.basictypes.Float64);
    $t.fastbox(3.14159, $g.____testlib.basictypes.Float64);
    $t.fastbox(20, $g.____testlib.basictypes.Integer);
    $t.fastbox(42, $g.____testlib.basictypes.Float64);
    $t.fastbox(42.5, $g.____testlib.basictypes.Float64);
    return;
  };
});
