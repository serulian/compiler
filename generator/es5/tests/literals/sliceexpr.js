$module('sliceexpr', function () {
  var $static = this;
  $static.TEST = function () {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Boolean).overArray([$t.box(false, $g.____testlib.basictypes.Boolean), $t.box(true, $g.____testlib.basictypes.Boolean), $t.box(false, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              return $result0.$index($t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $state.current = 1;
                $continue($state);
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
