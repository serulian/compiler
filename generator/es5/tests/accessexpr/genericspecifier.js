$module('genericspecifier', function () {
  var $instance = this;
  this.cls('SomeClass', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.$new = function () {
      var instance = new $static();
      function () {
      }.call(instance);
      return instance;
    };
  });
  $instance.SomeFunction = function (T, Q) {
    var $f = function () {
      return $promise.empty();
    };
    return $f;
  };
  $instance.DoSomething = function () {
    var $this = this;
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              SomeFunction($g.____graphs.typegraph.tests.testlib.basictypes.Integer, $g.____graphs.typegraph.tests.testlib.basictypes.String);
              SomeClass($g.____graphs.typegraph.tests.testlib.basictypes.Boolean);
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
