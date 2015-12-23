$module('continue', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
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
              4567;
              $state.current = 1;
              continue;

            case 3:
              2567;
              $state.current = -1;
              return;
          }
        }
      } catch (e) {
        $state.error = e;
        $state.current = -1;
        $callback($state);
      }
    };
    return $promise.build($state);
  };
});
