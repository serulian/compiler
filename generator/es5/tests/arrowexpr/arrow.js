$module('arrow', function () {
  var $static = this;
  $static.DoSomething = function (p) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var someint;
    $state.next = function ($callback) {
      try {
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
