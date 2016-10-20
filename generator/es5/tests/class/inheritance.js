$module('inheritance', function () {
  var $static = this;
  this.$class('16519abc', 'FirstClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeBool = $t.box(true, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    $instance.DoSomething = function () {
      var $this = this;
      return $promise.empty();
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "DoSomething|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('bcec6d2c', 'SecondClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeBool = $t.box(false, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    $instance.AnotherThing = function () {
      var $this = this;
      return $promise.empty();
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherThing|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('131aac76', 'MainClass', false, '', function () {
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
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "DoSomething|2|29dc432d<5ab5941e>": true,
        "AnotherThing|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
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
