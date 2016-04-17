$module('mini', function () {
  var $static = this;
  $static.TEST = function () {
    var lambda;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            lambda = function (someParam) {
              var $state = $t.sm(function ($continue) {
                $state.resolve($t.box(!$t.unbox(someParam), $g.____testlib.basictypes.Boolean));
                return;
              });
              return $promise.build($state);
            };
            lambda($t.box(false, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
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
