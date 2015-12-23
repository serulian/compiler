$module('nullcompare', function () {
  var $static = this;
  $static.DoSomething = function (someParam) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.returnValue = $op.nullcompare(someParam, 2);
            $state.current = -1;
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
