$module('numeric', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($continue) {
      $t.box(2, $g.____testlib.basictypes.Integer);
      $t.box(2.0, $g.____testlib.basictypes.Float64);
      $t.box(3.14159, $g.____testlib.basictypes.Float64);
      $t.box(20, $g.____testlib.basictypes.Integer);
      $t.box(42, $g.____testlib.basictypes.Float64);
      $t.box(42.5, $g.____testlib.basictypes.Float64);
      $state.resolve();
    });
    return $promise.build($state);
  };
});
