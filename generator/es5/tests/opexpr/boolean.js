$module('boolean', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            first = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
            second = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            $promise.resolve(first.$wrapped).then(function ($result2) {
              return $promise.resolve($result2 && second.$wrapped).then(function ($result1) {
                return $promise.resolve($result1 || first.$wrapped).then(function ($result0) {
                  $result = $t.fastbox($result0 || !second.$wrapped, $g.____testlib.basictypes.Boolean);
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
