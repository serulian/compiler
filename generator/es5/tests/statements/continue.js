$module('continue', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.box(1234, $g.____testlib.basictypes.Integer);
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
            $t.box(4567, $g.____testlib.basictypes.Integer);
            $state.current = 1;
            continue;

          case 3:
            $t.box(2567, $g.____testlib.basictypes.Integer);
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
