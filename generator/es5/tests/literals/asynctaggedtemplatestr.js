$module('asynctaggedtemplatestr', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('f7b645e4', function () {
    return $t.fastbox(30, $g.________testlib.basictypes.Integer);
  });
  $static.myFunction = $t.markpromising(function (pieces, values) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.asynctaggedtemplatestr.DoSomethingAsync()).then(function ($result0) {
              $result = $t.fastbox(($t.cast(values.$index($t.fastbox(0, $g.________testlib.basictypes.Integer)), $g.________testlib.basictypes.Integer, false).$wrapped + values.Length().$wrapped) + $result0.$wrapped, $g.________testlib.basictypes.Integer);
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
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            a = $t.fastbox(10, $g.________testlib.basictypes.Integer);
            b = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
            $promise.maybe($g.asynctaggedtemplatestr.myFunction($g.________testlib.basictypes.Slice($g.________testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.________testlib.basictypes.String), $t.fastbox("! It is ", $g.________testlib.basictypes.String), $t.fastbox("!", $g.________testlib.basictypes.String)]), $g.________testlib.basictypes.Slice($g.________testlib.basictypes.Stringable).overArray([a, b]))).then(function ($result0) {
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
            $resolve($t.fastbox(result.$wrapped == 42, $g.________testlib.basictypes.Boolean));
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
