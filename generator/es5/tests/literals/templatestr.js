$module('templatestr', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var result;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            a = $t.nominalwrap(1, $g.____testlib.basictypes.Integer);
            b = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.nominalwrap('This function is #', $g.____testlib.basictypes.String), $t.nominalwrap('! It is ', $g.____testlib.basictypes.String), $t.nominalwrap('!', $g.____testlib.basictypes.String)]).then(function ($result0) {
              return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]).then(function ($result1) {
                return $g.____testlib.basictypes.formatTemplateString($result0, $result1).then(function ($result2) {
                  $result = $result2;
                  $state.current = 1;
                  $callback($state);
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            result = $result;
            $g.____testlib.basictypes.String.$equals(result, $t.nominalwrap('This function is #1! It is true!', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
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
