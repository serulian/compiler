$module('genericspecifier', function () {
  var $static = this;
  $static.SomeFunction = function (T, Q) {
    var $f = function () {
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    };
    return $f;
  };
  $static.TEST = function () {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.genericspecifier.SomeFunction($g.____testlib.basictypes.Integer, $g.____testlib.basictypes.String)().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
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
