$module('nullcompare', function () {
  var $instance = this;
  $instance.DoSomething = function (someParam) {
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
              $state.returnValue = $op.nullcompare(someParam, 2);
              $state.current = -1;
              $callback($state);
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
