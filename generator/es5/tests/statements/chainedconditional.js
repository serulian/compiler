$module('chainedconditional', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            if (true) {
              $state.current = 1;
            } else {
              $state.current = 2;
            }
            continue;

          case 1:
            123;
            $state.current = 6;
            continue;

          case 2:
            if (false) {
              $state.current = 3;
            } else {
              $state.current = 4;
            }
            continue;

          case 3:
            456;
            $state.current = 5;
            continue;

          case 4:
            789;
            $state.current = 5;
            continue;

          case 5:
            $state.current = 6;
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
