$module('numeric', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.nominalwrap(2, $g.____testlib.basictypes.Integer);
            $t.nominalwrap(2.0, $g.____testlib.basictypes.Float64);
            $t.nominalwrap(3.14159, $g.____testlib.basictypes.Float64);
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
