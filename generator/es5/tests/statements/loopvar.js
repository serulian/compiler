$module('loopvar', function () {
  var $static = this;
  $static.DoSomething = function (somethingElse) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            1234;
            $state.current = 1;
            continue;

          case 1:
            if (somethingElse) {
              $state.current = 2;
              continue;
            } else {
              $state.current = 3;
              continue;
            }
            break;

          case 2:
            7654;
            $state.current = 1;
            continue;

          case 3:
            5678;
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
