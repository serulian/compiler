$module('conditionalelse', function () {
  var $instance = this;
  $instance.DoSomething = function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
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
