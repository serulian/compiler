$module('full', function () {
  var $static = this;
  $static.TEST = $t.markpromising(function () {
    var $result;
    var lambda;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            lambda = function (firstParam, secondParam) {
              return secondParam;
            };
            $promise.maybe(lambda($t.fastbox(123, $g.________testlib.basictypes.Integer), $t.fastbox(true, $g.________testlib.basictypes.Boolean))).then(function ($result0) {
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
