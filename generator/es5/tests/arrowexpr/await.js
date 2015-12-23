$module('await', function () {
  var $static = this;
  $static.DoSomething = function (p) {
    var $awaitresult$2;
    var $state = $t.sm(function ($callback) {
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

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
