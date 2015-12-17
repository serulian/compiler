$module('await', function () {
  var $instance = this;
  $instance.DoSomething = function (p) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var $awaitresult$2;
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              p.then(function (returnValue) {
                $state.current = 1;
                $awaitresult$2 = returnValue;
                $state.next($callback);
              }).catch(function (e) {
                $state.error = e;
                $state.current = -1;
                $callback($state);
              });
              return;

            case 1:
              $state.returnValue = $awaitresult$2;
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
