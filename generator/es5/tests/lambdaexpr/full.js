$module('full', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            function (firstParam, secondParam) {
              var $state = $t.sm(function ($callback) {
                while (true) {
                  switch ($state.current) {
                    case 0:
                      1234;
                      $state.returnValue = 4567;
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
