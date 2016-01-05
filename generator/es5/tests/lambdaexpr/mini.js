$module('mini', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            function (someParam) {
              var $state = $t.sm(function ($callback) {
                while (true) {
                  switch ($state.current) {
                    case 0:
                      $state.returnValue = someParam;
                      $state.current = -1;
                      $state.returnValue = null;
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
            $state.returnValue = null;
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
