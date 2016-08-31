$module('map', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var map;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Map($g.____testlib.basictypes.String, $g.____testlib.basictypes.Boolean).forArrays([$t.box('hello', $g.____testlib.basictypes.String), $t.box('hi', $g.____testlib.basictypes.String)], [$t.box(true, $g.____testlib.basictypes.Boolean), $t.box(false, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            map = $result;
            map.$index($t.box('hello', $g.____testlib.basictypes.String)).then(function ($result2) {
              return $promise.resolve($result2).then(function ($result1) {
                return $promise.resolve($t.unbox($t.nullcompare($result1, $t.box(false, $g.____testlib.basictypes.Boolean)))).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || map.$index($t.box('hi', $g.____testlib.basictypes.String))).then(function ($result4) {
                    return $promise.resolve($result4).then(function ($result3) {
                      $result = $t.box($result0 && !$t.unbox($t.nullcompare($result3, $t.box(true, $g.____testlib.basictypes.Boolean))), $g.____testlib.basictypes.Boolean);
                      $current = 2;
                      $continue($resolve, $reject);
                      return;
                    });
                  });
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
