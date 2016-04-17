$module('taggedtemplatestr', function () {
  var $static = this;
  $static.myFunction = function (pieces, values) {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            values.$index($t.box(0, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              return values.Length().then(function ($result1) {
                return $g.____testlib.basictypes.Integer.$plus($t.cast($result0, $g.____testlib.basictypes.Integer), $result1).then(function ($result2) {
                  $result = $result2;
                  $state.current = 1;
                  $continue($state);
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
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            a = $t.box(10, $g.____testlib.basictypes.Integer);
            b = $t.box(true, $g.____testlib.basictypes.Boolean);
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.box("This function is #", $g.____testlib.basictypes.String), $t.box("! It is ", $g.____testlib.basictypes.String), $t.box("!", $g.____testlib.basictypes.String)]).then(function ($result0) {
              return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]).then(function ($result1) {
                return $g.taggedtemplatestr.myFunction($result0, $result1).then(function ($result2) {
                  $result = $result2;
                  $state.current = 1;
                  $continue($state);
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            result = $result;
            $g.____testlib.basictypes.Integer.$equals(result, $t.box(12, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
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
