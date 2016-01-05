$module('loopexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            1234;
            $state.current = 1;
            continue;

          case 1:
            if (true) {
              $state.current = 2;
            } else {
              $state.current = 3;
            }
            continue;

          case 2:
            1357;
            $state.current = 1;
            continue;

          case 3:
            5678;
            $state.current = -1;
            $state.returnValue = null;
            $callback($state);
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
