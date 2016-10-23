$module('functioncall', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.functioncall.AnotherFunction($t.fastbox(false, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
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
  $static.AnotherFunction = function (param) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.fastbox(!param.$wrapped, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
});
