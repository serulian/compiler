$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.returnValue = 2;
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
  this.$init(function () {
    return $promise.wrap(function () {
      $this.someInt = 2;
    });
  });
  this.$init(function () {
    var $returnValue$1;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.basic.AnotherFunction().then(function (returnValue) {
              $state.current = 1;
              $returnValue$1 = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 1:
            $this.anotherInt = $returnValue$1;
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
  });
});
