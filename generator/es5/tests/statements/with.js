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
            $t.pushr($state, '$withExpr$1', $withExpr$1);
            456;
            $t.popr($state, '$withExpr$1');
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
