$module('with', function () {
  var $instance = this;
  $instance.DoSomething = function (someExpr) {
    var $this = this;
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
              $state.$pushRAII('$withExpr$1', $withExpr$1);
              456;
              $state.$popRAII('$withExpr$1');
              789;
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
