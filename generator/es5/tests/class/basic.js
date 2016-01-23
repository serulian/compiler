$module('basic', function () {
  var $static = this;
  this.$class('SomeClass', false, function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve(2).then(function (result) {
        instance.SomeInt = result;
      }));
      init.push($g.basic.CoolFunction().then(function ($result0) {
        return $promise.resolve($result0);
      }).then(function (result) {
        instance.AnotherBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.AnotherFunction = function () {
      var $this = this;
      return $promise.empty();
    };
  });

  $static.CoolFunction = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.resolve(true);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.basic.SomeClass.new().then(function ($result0) {
              $result = $result0.AnotherBool;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $state.resolve($result);
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
