$module('taggedtemplatestr', function () {
  var $static = this;
  $static.myFunction = function (pieces, values) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            values.$index($t.fastbox(0, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              return values.Length().then(function ($result1) {
                $result = $t.fastbox($t.cast($result0, $g.____testlib.basictypes.Integer, false).$wrapped + $result1.$wrapped, $g.____testlib.basictypes.Boolean);
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
  $static.TEST = function () {
    var $result;
    var a;
    var b;
    var result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            a = $t.fastbox(10, $g.____testlib.basictypes.Integer);
            b = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.____testlib.basictypes.String), $t.fastbox("! It is ", $g.____testlib.basictypes.String), $t.fastbox("!", $g.____testlib.basictypes.String)]).then(function ($result1) {
              return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]).then(function ($result2) {
                return $g.taggedtemplatestr.myFunction($result1, $result2).then(function ($result0) {
                  $result = $result0;
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
            result = $result;
            $resolve($t.fastbox(result.$wrapped == 12, $g.____testlib.basictypes.Boolean));
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
