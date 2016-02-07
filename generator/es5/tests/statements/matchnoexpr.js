$module('matchnoexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.nominalwrap(123, $g.____testlib.basictypes.Integer);
            if (false) {
              $state.current = 1;
              continue;
            } else {
              $state.current = 3;
              continue;
            }
            break;

          case 1:
            $t.nominalwrap(1234, $g.____testlib.basictypes.Integer);
            $state.current = 2;
            continue;

          case 2:
            $t.nominalwrap(789, $g.____testlib.basictypes.Integer);
            $state.current = -1;
            return;

          case 3:
            if (true) {
              $state.current = 4;
              continue;
            } else {
              $state.current = 5;
              continue;
            }
            break;

          case 4:
            $t.nominalwrap(2345, $g.____testlib.basictypes.Integer);
            $state.current = 2;
            continue;

          case 5:
            if (true) {
              $state.current = 6;
              continue;
            } else {
              $state.current = 2;
              continue;
            }
            break;

          case 6:
            $t.nominalwrap(3456, $g.____testlib.basictypes.Integer);
            $state.current = 2;
            continue;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
