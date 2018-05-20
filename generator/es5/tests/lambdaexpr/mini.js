$module('mini', function () {
  var $static = this;
  $static.TEST = $t.markpromising(function () {
    var $result;
    var lambda;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            lambda = function (someParam) {
              return $t.fastbox(!someParam.$wrapped, $g.________testlib.basictypes.Boolean);
            };
            $promise.maybe(lambda($t.fastbox(false, $g.________testlib.basictypes.Boolean))).then(function ($result0) {
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
});
