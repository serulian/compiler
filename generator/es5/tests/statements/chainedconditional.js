$module('chainedconditional', function () {
  var $static = this;
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            if (false) {
              $state.current = 1;
              continue;
            } else {
              $state.current = 2;
              continue;
            }
            break;

          case 1:
            $t.nominalwrap(123, $g.____testlib.basictypes.Integer);
            $state.resolve($t.nominalwrap(false, $g.____testlib.basictypes.Boolean));
            return;

          case 2:
            if (false) {
              $state.current = 3;
              continue;
            } else {
              $state.current = 4;
              continue;
            }
            break;

          case 3:
            $t.nominalwrap(456, $g.____testlib.basictypes.Integer);
            $state.resolve($t.nominalwrap(false, $g.____testlib.basictypes.Boolean));
            return;

          case 4:
            $t.nominalwrap(789, $g.____testlib.basictypes.Integer);
            $state.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
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
