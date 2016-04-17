$module('boolean', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            first = $t.box(true, $g.____testlib.basictypes.Boolean);
            second = $t.box(false, $g.____testlib.basictypes.Boolean);
            $promise.resolve($t.unbox(first)).then(function ($result2) {
              return $promise.resolve($result2 && $t.unbox(second)).then(function ($result1) {
                return $promise.resolve($result1 || $t.unbox(first)).then(function ($result0) {
                  $result = $t.box($result0 || !$t.unbox(second), $g.____testlib.basictypes.Boolean);
                  $state.current = 1;
                  $continue($state);
                });
              });
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
