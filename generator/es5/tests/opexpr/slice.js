$module('slice', function () {
  var $static = this;
  this.$class('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.DoSomething = function (c) {
    var $returnValue$1;
    var $returnValue$2;
    var $returnValue$3;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            c.$slice(1, 2).then(function (returnValue) {
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
            c.$slice(null, 1).then(function (returnValue) {
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
            c.$slice(1, null).then(function (returnValue) {
              $state.current = 3;
              $returnValue$3 = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 3:
            $returnValue$3;
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
