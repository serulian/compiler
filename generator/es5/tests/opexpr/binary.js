$module('binary', function () {
  var $instance = this;
  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.$new = function () {
      var instance = new $static();
      function () {
      }.call(instance);
      return instance;
    };
  });
  $instance.DoSomething = function (first, second) {
    var $this = this;
    var $state = {
      current: 0,
      returnValue: null,
    };
    var $returnValue$1;
    var $returnValue$2;
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              plus(first, second).then(function (returnValue) {
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
              $returnValue$1;
              minus(first, second).then(function (returnValue) {
                $state.current = 2;
                $returnValue$2 = returnValue;
                $state.next($callback);
              }).catch(function (e) {
                $state.error = e;
                $state.current = -1;
                $callback($state);
              });
              return;

            case 2:
              $returnValue$2;
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
