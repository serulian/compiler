$module('functioncall', function () {
  var $static = this;
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.functioncall.AnotherFunction($t.nominalwrap(false, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
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
  $static.AnotherFunction = function (param) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.resolve($t.nominalwrap(!$t.nominalunwrap(param), $g.____testlib.basictypes.Boolean));
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
