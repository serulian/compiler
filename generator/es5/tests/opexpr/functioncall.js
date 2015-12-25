$module('functioncall', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $returnValue$1;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.functioncall.AnotherFunction(2).then(function (returnValue) {
              $state.current = 1;
              $returnValue$1 = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 1:
            $returnValue$1;
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
  $static.AnotherFunction = function (someparam) {
    return $promise.empty();
  };
});
