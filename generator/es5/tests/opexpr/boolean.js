$module('boolean', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            first = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
            second = $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
            $promise.resolve($t.nominalunwrap(first)).then(function ($result2) {
              return $promise.resolve($result2 && $t.nominalunwrap(second)).then(function ($result1) {
                return $promise.resolve($result1 || $t.nominalunwrap(first)).then(function ($result0) {
                  $result = $t.nominalwrap($result0 || !$t.nominalunwrap(second), $g.____testlib.basictypes.Boolean);
                  $state.current = 1;
                  $callback($state);
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
