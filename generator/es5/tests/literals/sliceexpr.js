$module('sliceexpr', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Boolean).overArray([$t.fastbox(false, $g.____testlib.basictypes.Boolean), $t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox(false, $g.____testlib.basictypes.Boolean)]).then(function ($result1) {
              return $result1.$index($t.fastbox(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
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
