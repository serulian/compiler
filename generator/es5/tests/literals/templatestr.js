$module('templatestr', function () {
  var $static = this;
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
            a = $t.box(1, $g.____testlib.basictypes.Integer);
            b = $t.box(true, $g.____testlib.basictypes.Boolean);
            $g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.box("This function is #", $g.____testlib.basictypes.String), $t.box("! It is ", $g.____testlib.basictypes.String), $t.box("!", $g.____testlib.basictypes.String)]).then(function ($result1) {
              return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]).then(function ($result2) {
                return $g.____testlib.basictypes.formatTemplateString($result1, $result2).then(function ($result0) {
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
            $g.____testlib.basictypes.String.$equals(result, $t.box('This function is #1! It is true!', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
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
