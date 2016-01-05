$module('structuralcast', function () {
  var $static = this;
  this.$class('BaseClass', function (T) {
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

  this.$class('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push(function () {
        var $this = this;
        return $g.structuralcast.BaseClass($g.____graphs.srg.typeconstructor.tests.testlib.basictypes.Integer).new().then(function (value) {
          $this.BaseClass$Integer = value;
        });
      }());
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.DoSomething = function (sc) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            sc.BaseClass$Integer;
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
