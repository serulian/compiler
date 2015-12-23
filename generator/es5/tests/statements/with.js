$module('with', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var $withExpr$1;
    $state.next = function ($callback) {
      try {
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
      } catch (e) {
        $state.error = e;
        $state.current = -1;
        $callback($state);
      }
    };
    return $promise.build($state);
  };
});
