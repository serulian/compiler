$module('boolean', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            first = $t.box(true, $g.____testlib.basictypes.Boolean);
            second = $t.box(false, $g.____testlib.basictypes.Boolean);
            $promise.resolve($t.unbox(first)).then(function ($result2) {
              return $promise.resolve($result2 && $t.unbox(second)).then(function ($result1) {
                return $promise.resolve($result1 || $t.unbox(first)).then(function ($result0) {
                  $result = $t.box($result0 || !$t.unbox(second), $g.____testlib.basictypes.Boolean);
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
