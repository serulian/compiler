$module('withas', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var someName;
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              123;
              someName = someExpr;
              $state.$pushRAII('someName', someName);
              456;
              $state.$popRAII('someName');
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
