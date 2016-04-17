$module('mixed', function () {
  var $static = this;
  $static.TEST = function () {
    var finalIndex;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            finalIndex = $t.box(-2, $g.____testlib.basictypes.Integer);
            $g.____testlib.basictypes.Integer.$compare(finalIndex, $t.box(10, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($t.unbox($result1) >= 0).then(function ($result0) {
                return ($promise.shortcircuit($result0, true) || $g.____testlib.basictypes.Integer.$compare(finalIndex, $t.box(0, $g.____testlib.basictypes.Integer))).then(function ($result2) {
                  $result = $t.box($result0 || ($t.unbox($result2) < 0), $g.____testlib.basictypes.Boolean);
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
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
