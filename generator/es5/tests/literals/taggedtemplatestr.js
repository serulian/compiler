$module('taggedtemplatestr', function () {
  var $static = this;
  $static.myFunction = function (pieces, values) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            values.$index($t.nominalwrap(0, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              return values.Length().then(function ($result1) {
                return $g.____testlib.basictypes.Integer.$plus($t.cast($result0, $g.____testlib.basictypes.Integer), $result1).then(function ($result2) {
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
  $static.TEST = function () {
    var a;
    var b;
    var result;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            a = $t.nominalwrap(10, $g.____testlib.basictypes.Integer);
            b = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.nominalwrap('This function is #', $g.____testlib.basictypes.String), $t.nominalwrap('! It is ', $g.____testlib.basictypes.String), $t.nominalwrap('!', $g.____testlib.basictypes.String)]).then(function ($result0) {
              return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]).then(function ($result1) {
                return $g.taggedtemplatestr.myFunction($result0, $result1).then(function ($result2) {
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
            $g.____testlib.basictypes.Integer.$equals(result, $t.nominalwrap(12, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
