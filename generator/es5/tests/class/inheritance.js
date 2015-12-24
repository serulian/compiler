$module('inheritance', function () {
  var $static = this;
  this.$class('FirstClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push(function () {
        return $promise.wrap(function () {
          $this.SomeInt = 2;
        });
      }());
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function () {
      return $promise.empty();
    };
  });

  this.$class('SecondClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push(function () {
        return $promise.wrap(function () {
          $this.SomeInt = 4;
        });
      }());
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.AnotherThing = function () {
      return $promise.empty();
    };
  });

  this.$class('MainClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push(function () {
        var $this = this;
        return $g.inheritance.FirstClass.new().then(function (value) {
          $this.FirstClass = value;
        });
      }());
      init.push(function () {
        var $this = this;
        return $g.inheritance.SecondClass.new().then(function (value) {
          $this.SecondClass = value;
        });
      }());
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function (i) {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.returnValue = $this.SomeInt;
              $state.current = -1;
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
    Object.defineProperty($instance, 'SomeInt', {
      get: function () {
        return $instance.FirstClass.SomeInt;
      },
      set: function (val) {
        $instance.FirstClass.SomeInt = val;
      },
    });
    Object.defineProperty($instance, 'AnotherThing', {
      get: function () {
        return $instance.SecondClass.AnotherThing;
      },
    });
  });

});
