$module('break', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.nominalwrap(1234, $g.____testlib.basictypes.Integer);
            $state.current = 1;
            continue;

          case 1:
            if (true) {
              $state.current = 2;
              continue;
            } else {
              $state.current = 3;
              continue;
            }
            break;

          case 2:
            $t.nominalwrap(4567, $g.____testlib.basictypes.Integer);
            $state.current = 3;
            continue;

          case 3:
            $t.nominalwrap(2567, $g.____testlib.basictypes.Integer);
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
