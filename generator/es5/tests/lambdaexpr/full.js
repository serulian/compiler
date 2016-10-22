$module('full', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var lambda;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            lambda = function (firstParam, secondParam) {
              var $current = 0;
              var $continue = function ($resolve, $reject) {
                $resolve(secondParam);
                return;
              };
              return $promise.new($continue);
            };
            lambda($t.fastbox(123, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
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
  };
});
