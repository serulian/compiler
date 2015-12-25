$module('conditionalelse', function () {
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
            $state.current = 3;
            continue;

          case 2:
            456;
            $state.current = 3;
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
