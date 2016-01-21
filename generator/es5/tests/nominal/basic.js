$module('basic', function () {
  var $static = this;
  this.$class('SomeClass', false, function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function () {
      var $this = this;
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
  });

  this.$type('MyType', function () {
    var $instance = this;
    $instance.AnotherThing = function ($this) {
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $this.DoSomething().then(function ($result0) {
                $result = $result0;
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
    $instance.SomeProp = $t.property(true, function ($this) {
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
    });
  });

  $static.TEST = function () {
    var m;
    var sc;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.basic.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            m = sc;
            $g.basic.MyType.SomeProp(m).then(function ($result0) {
              return $g.basic.MyType.AnotherThing(m).then(function ($result1) {
                $result = $result0 && $result1;
                $state.current = 2;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
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
