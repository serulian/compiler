$module('matchexpr', function () {
  var $instance = this;
  $instance.DoSomething = function (someVar) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              123;
              if (1 != someVar) {
                $state.current = 1;
                continue;
              }
              1234;
              $state.current = 3;
              continue;

            case 1:
              if (2 != someVar) {
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
