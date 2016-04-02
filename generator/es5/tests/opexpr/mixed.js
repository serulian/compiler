$module('mixed', function () {
  var $static = this;
  $static.TEST = function () {
    var finalIndex;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            finalIndex = $t.nominalwrap(-2, $g.____testlib.basictypes.Integer);
            $g.____testlib.basictypes.Integer.$compare(finalIndex, $t.nominalwrap(10, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($t.nominalunwrap($result1) >= 0).then(function ($result0) {
                return ($promise.shortcircuit($result0, true) || $g.____testlib.basictypes.Integer.$compare(finalIndex, $t.nominalwrap(0, $g.____testlib.basictypes.Integer))).then(function ($result2) {
                  $result = $t.nominalwrap($result0 || ($t.nominalunwrap($result2) < 0), $g.____testlib.basictypes.Boolean);
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
