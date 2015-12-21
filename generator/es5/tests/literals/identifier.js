$module('identifier', function () {
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
  $instance.DoSomething = function (someParam) {
    var $this = this;
    var $state = {
      current: 0,
      returnValue: null,
    };
    var someVar;
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              someVar = 2;
              $g.identifier.SomeClass;
              $g.identifier.AnotherFunction;
              someParam;
              someVar;
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
  $instance.AnotherFunction = function () {
    return $promise.empty();
  };
});
