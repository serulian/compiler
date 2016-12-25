$module('async', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('088145a7', function () {
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  });
  $static.DoSomethingElse = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.async.DoSomethingAsync()).then(function ($result0) {
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
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var a;
    var b;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.async.DoSomethingElse()).then(function ($result0) {
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
            try {
              var $expr = $result;
              a = $expr;
              b = null;
            } catch ($rejected) {
              b = $rejected;
              a = null;
            }
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 2:
            $resolve(a);
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
