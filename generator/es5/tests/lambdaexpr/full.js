$module('full', function () {
  var $static = this;
  $static.TEST = function () {
    var lambda;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            lambda = function (firstParam, secondParam) {
              var $state = $t.sm(function ($callback) {
                while (true) {
                  switch ($state.current) {
                    case 0:
                      $state.resolve(secondParam);
                      return;

                    default:
                      $state.current = -1;
                      return;
                  }
                }
              });
              return $promise.build($state);
            };
            lambda($t.box(123, $g.____testlib.basictypes.Integer), $t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $state.resolve($result);
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
