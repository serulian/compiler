$module('indexer', function () {
  var $static = this;
  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function ($callback) {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.DoSomething = function (c) {
    var $returnValue$1;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            c.$index(true).then(function (returnValue) {
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
