$module('async', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('9a368168', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.async.DoSomethingAsync()).then(function ($result1) {
              return $promise.resolve($result1.$wrapped).then(function ($result0) {
                $result = $result0 ? $t.fastbox(true, $g.________testlib.basictypes.Boolean) : $t.fastbox(false, $g.________testlib.basictypes.Boolean);
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
  });
});
