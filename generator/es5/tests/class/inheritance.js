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
    this.$typesig = function () {
      return $t.createtypesig(['DoSomething', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.inheritance.FirstClass).$typeref()]);
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
    this.$typesig = function () {
      return $t.createtypesig(['AnotherThing', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.inheritance.SecondClass).$typeref()]);
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
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($this.SomeBool);
        return;
      };
      return $promise.new($continue);
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
    this.$typesig = function () {
      return $t.createtypesig(['DoSomething', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['AnotherThing', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.inheritance.MainClass).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.inheritance.MainClass.new().then(function ($result1) {
              return $result1.DoSomething().then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
