$module('sliceexpr', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Boolean).overArray([$t.box(false, $g.____testlib.basictypes.Boolean), $t.box(true, $g.____testlib.basictypes.Boolean), $t.box(false, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              return $result0.$index($t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
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
