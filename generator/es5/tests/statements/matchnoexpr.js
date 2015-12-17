$module('matchnoexpr', function () {
  var $instance = this;
  $instance.DoSomething = function () {
    var $this = this;
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
