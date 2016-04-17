$module('numeric', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.box(2, $g.____testlib.basictypes.Integer);
      $t.box(2.0, $g.____testlib.basictypes.Float64);
      $t.box(3.14159, $g.____testlib.basictypes.Float64);
      $t.box(20, $g.____testlib.basictypes.Integer);
      $t.box(42, $g.____testlib.basictypes.Float64);
      $t.box(42.5, $g.____testlib.basictypes.Float64);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
