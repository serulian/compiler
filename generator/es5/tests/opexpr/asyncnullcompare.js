$module('asyncnullcompare', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('9be4de44', function () {
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var someBool;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            someBool = null;
            $promise.resolve(someBool).then(function ($result0) {
              return ($promise.shortcircuit($result0, null) || $promise.translate($g.asyncnullcompare.DoSomethingAsync())).then(function ($result1) {
                $result = $t.asyncnullcompare($result0, $result1);
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
