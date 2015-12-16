$module('varassign', function () {
  var $instance = this;
  $instance.DoSomething = function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var i;
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              i = 2;
              i = 3;
              1234;
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
