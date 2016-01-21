$module('conditionalelse', function () {
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
            $state.resolve(false);
            return;

          case 2:
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
