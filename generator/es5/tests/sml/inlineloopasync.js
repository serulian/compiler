$module('inlineloopasync', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('a3e83ce9', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.somenumber = $t.markpromising(function (props, value) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.inlineloopasync.DoSomethingAsync()).then(function ($result0) {
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
            $resolve(value);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.somesum = function (props, numbers) {
    return $t.fastbox(42, $g.________testlib.basictypes.Integer);
  };
  $static.TEST = $t.markpromising(function () {
    var result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      result = $g.inlineloopasync.somesum($g.________testlib.basictypes.Mapping($t.any).Empty(), $g.________testlib.basictypes.MapStream($g.________testlib.basictypes.Integer, $g.________testlib.basictypes.Integer)($g.________testlib.basictypes.Integer.$range($t.fastbox(0, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer)), $t.markpromising(function (index) {
        var $result;
        var $current = 0;
        var $continue = function ($resolve, $reject) {
          localasyncloop: while (true) {
            switch ($current) {
              case 0:
                $promise.maybe($g.inlineloopasync.somenumber($g.________testlib.basictypes.Mapping($t.any).Empty(), index)).then(function ($result0) {
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
                $resolve($result);
                return;

              default:
                $resolve();
                return;
            }
          }
        };
        return $promise.new($continue);
      })));
      $resolve($t.fastbox(result.$wrapped == 42, $g.________testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  });
});
