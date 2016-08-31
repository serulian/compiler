$module('withas', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var $result;
    var someName;
    var $current = 0;
    var $resources = $t.resourcehandler();
    var $continue = function ($resolve, $reject) {
      $resolve = $resources.bind($resolve);
      $reject = $resources.bind($reject);
      while (true) {
        switch ($current) {
          case 0:
            $t.box(123, $g.____testlib.basictypes.Integer);
            someName = someExpr;
            $resources.pushr(someName, 'someName');
            $t.box(456, $g.____testlib.basictypes.Integer);
            $resources.popr('someName').then(function ($result0) {
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
            $t.box(789, $g.____testlib.basictypes.Integer);
            $resolve();
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
