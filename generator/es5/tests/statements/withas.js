$module('withas', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var someName;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            123;
            someName = someExpr;
            $state.pushr(someName, 'someName');
            456;
            $state.popr('someName').then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
            789;
            $state.current = -1;
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
