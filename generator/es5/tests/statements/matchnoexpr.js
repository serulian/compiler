$module('matchnoexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            123;
            if (false) {
              $state.current = 1;
              continue;
            } else {
              $state.current = 3;
              continue;
            }
            break;

          case 1:
            1234;
            $state.current = 2;
            continue;

          case 2:
            789;
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
            2345;
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
            3456;
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
