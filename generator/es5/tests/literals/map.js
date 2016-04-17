$module('map', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.Map($g.____testlib.basictypes.String, $g.____testlib.basictypes.Boolean).forArrays([$t.box('hello', $g.____testlib.basictypes.String), $t.box('hi', $g.____testlib.basictypes.String)], [$t.box(true, $g.____testlib.basictypes.Boolean), $t.box(false, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            map = $result;
            map.$index($t.box('hello', $g.____testlib.basictypes.String)).then(function ($result1) {
              return $promise.resolve($t.unbox($t.nullcompare($result1, $t.box(false, $g.____testlib.basictypes.Boolean)))).then(function ($result0) {
                return ($promise.shortcircuit($result0, false) || map.$index($t.box('hi', $g.____testlib.basictypes.String))).then(function ($result2) {
                  $result = $t.box($result0 && !$t.unbox($t.nullcompare($result2, $t.box(true, $g.____testlib.basictypes.Boolean))), $g.____testlib.basictypes.Boolean);
                  $state.current = 2;
                  $continue($state);
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
