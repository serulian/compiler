$module('list', function () {
  var $static = this;
  $static.TEST = function () {
    var l;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.List($t.any).forArray([$t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer), $t.box(3, $g.____testlib.basictypes.Integer), $t.box(true, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            l = $result;
            l.Count().then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.box(4, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $state.current = 2;
                $callback($state);
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
