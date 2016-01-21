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
            123;
            $state.resolve(false);
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
            456;
            $state.resolve(false);
            return;

          case 4:
            789;
            $state.resolve(true);
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
