$module('arrow', function () {
  var $static = this;
  $static.DoSomething = function (p) {
    var someint;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            p.then(function (returnValue) {
              $state.current = 1;
              someint = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 1:
            someint;
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
