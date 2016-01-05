$module('with', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var $withExpr$1;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            123;
            $withExpr$1 = someExpr;
            $state.pushr('$withExpr$1', $withExpr$1);
            456;
            $state.popr('$withExpr$1');
            789;
            $state.current = -1;
            $state.returnValue = null;
            $callback($state);
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
