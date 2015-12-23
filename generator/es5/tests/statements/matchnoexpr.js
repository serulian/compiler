$module('matchnoexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            123;
            if (true != true) {
              $state.current = 1;
              continue;
            }
            1234;
            $state.current = 3;
            continue;

          case 1:
            if (false != true) {
              $state.current = 2;
              continue;
            }
            2345;
            $state.current = 3;
            continue;

          case 2:
            3456;
            $state.current = 3;
            continue;

          case 3:
            789;
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
