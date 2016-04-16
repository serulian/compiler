$module('inheritance', function () {
  var $static = this;
  this.$class('FirstClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.SomeBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function () {
      var $this = this;
      return $promise.empty();
    };
  });

  this.$class('SecondClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.SomeBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.AnotherThing = function () {
      var $this = this;
      return $promise.empty();
    };
  });

  this.$class('MainClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($g.inheritance.FirstClass.new().then(function (value) {
        instance.FirstClass = value;
      }));
      init.push($g.inheritance.SecondClass.new().then(function (value) {
        instance.SecondClass = value;
      }));
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
              $state.resolve($this.SomeBool);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
    Object.defineProperty($instance, 'SomeBool', {
      get: function () {
        return this.FirstClass.SomeBool;
      },
      set: function (val) {
        this.FirstClass.SomeBool = val;
      },
    });
    Object.defineProperty($instance, 'AnotherThing', {
      get: function () {
        return this.SecondClass.AnotherThing;
      },
    });
  });

  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.inheritance.MainClass.new().then(function ($result0) {
              return $result0.DoSomething().then(function ($result1) {
                $result = $result1;
                $state.current = 1;
                $callback($state);
              });
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
