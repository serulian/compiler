$module('genericspecifier', function () {
  var $static = this;
  this.$class('SomeClass', function (T) {
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

  $static.SomeFunction = function (T, Q) {
    var $f = function () {
      return $promise.empty();
    };
    return $f;
  };
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.genericspecifier.SomeFunction($g.____graphs.srg.typeconstructor.tests.testlib.basictypes.Integer, $g.____graphs.srg.typeconstructor.tests.testlib.basictypes.String);
            $g.genericspecifier.SomeClass($g.____graphs.srg.typeconstructor.tests.testlib.basictypes.Boolean);
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
