$module('varassign', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            i = $t.nominalwrap(2, $g.____testlib.basictypes.Integer);
            i = $t.nominalwrap(3, $g.____testlib.basictypes.Integer);
            $t.nominalwrap(1234, $g.____testlib.basictypes.Integer);
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
