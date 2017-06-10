$module('asynctaggedtemplatestr', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('9c7a98cc', function () {
    return $t.fastbox(30, $g.____testlib.basictypes.Integer);
  });
  $static.myFunction = $t.markpromising(function (pieces, values) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.asynctaggedtemplatestr.DoSomethingAsync()).then(function ($result0) {
              $result = $t.fastbox(($t.cast(values.$index($t.fastbox(0, $g.____testlib.basictypes.Integer)), $g.____testlib.basictypes.Integer, false).$wrapped + values.Length().$wrapped) + $result0.$wrapped, $g.____testlib.basictypes.Integer);
              $current = 1;
              $continue($resolve, $reject);
              return;
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
  });
  $static.TEST = $t.markpromising(function () {
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
            $promise.maybe($g.asynctaggedtemplatestr.myFunction($g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.____testlib.basictypes.String), $t.fastbox("! It is ", $g.____testlib.basictypes.String), $t.fastbox("!", $g.____testlib.basictypes.String)]), $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]))).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            result = $result;
            $resolve($t.fastbox(result.$wrapped == 42, $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
