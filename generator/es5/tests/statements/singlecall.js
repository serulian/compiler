$module('singlecall', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.fastbox(42, $g.____testlib.basictypes.Integer));
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
            $g.singlecall.DoSomething().then(function ($result1) {
              return $g.____testlib.basictypes.Integer.$equals($result1, $t.fastbox(42, $g.____testlib.basictypes.Integer)).then(function ($result0) {
                $result = $result0;
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
  };
});
