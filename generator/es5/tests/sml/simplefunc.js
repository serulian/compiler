$module('simplefunc', function () {
  var $static = this;
  $static.SimpleFunction = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.simplefunc.SimpleFunction().then(function ($result0) {
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
