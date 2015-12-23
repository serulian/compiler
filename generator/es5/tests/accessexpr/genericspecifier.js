$module('genericspecifier', function () {
  var $static = this;
  this.cls('SomeClass', function (T) {
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

  $static.SomeFunction = function (T, Q) {
    var $f = function () {
      return $promise.empty();
    };
    return $f;
  };
  $static.DoSomething = function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.genericspecifier.SomeFunction($g.____graphs.typegraph.tests.testlib.basictypes.Integer, $g.____graphs.typegraph.tests.testlib.basictypes.String);
              $g.genericspecifier.SomeClass($g.____graphs.typegraph.tests.testlib.basictypes.Boolean);
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
