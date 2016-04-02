$module('map', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.Map($g.____testlib.basictypes.String, $g.____testlib.basictypes.Boolean).forArrays([$t.nominalwrap('hello', $g.____testlib.basictypes.String), $t.nominalwrap('hi', $g.____testlib.basictypes.String)], [$t.nominalwrap(true, $g.____testlib.basictypes.Boolean), $t.nominalwrap(false, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            map = $result;
            map.$index($t.nominalwrap('hello', $g.____testlib.basictypes.String)).then(function ($result1) {
              return $promise.resolve($t.nominalunwrap($t.nullcompare($result1, $t.nominalwrap(false, $g.____testlib.basictypes.Boolean)))).then(function ($result0) {
                return ($promise.shortcircuit($result0, false) || map.$index($t.nominalwrap('hi', $g.____testlib.basictypes.String))).then(function ($result2) {
                  $result = $t.nominalwrap($result0 && !$t.nominalunwrap($t.nullcompare($result2, $t.nominalwrap(true, $g.____testlib.basictypes.Boolean))), $g.____testlib.basictypes.Boolean);
                  $state.current = 2;
                  $callback($state);
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
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
