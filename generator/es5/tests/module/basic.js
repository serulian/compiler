$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.returnValue = 2;
              $state.current = -1;
              $callback($state);
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
  this.$init(function () {
    return $promise.wrap(function () {
      $this.someInt = 2;
    });
  });
  this.$init(function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var $returnValue$1;
    $state.next = function ($callback) {
      try {
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
              $state.returnValue = $returnValue$1;
              $state.current = -1;
              $callback($state);
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
  });
});
