$module('numeric', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.box(2, $g.____testlib.basictypes.Integer);
            $t.box(2.0, $g.____testlib.basictypes.Float64);
            $t.box(3.14159, $g.____testlib.basictypes.Float64);
            $t.box(20, $g.____testlib.basictypes.Integer);
            $t.box(42, $g.____testlib.basictypes.Float64);
            $t.box(42.5, $g.____testlib.basictypes.Float64);
            $state.current = -1;
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
