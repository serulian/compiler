$module('asyncinheritance', function () {
  var $static = this;
  this.$class('fba41d60', 'FirstClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.translate($g.asyncinheritance.DoSomethingAsync()).then(function ($result0) {
        instance.SomeBool = $result0;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function () {
      var $this = this;
      return;
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

  this.$class('ce8b287b', 'SecondClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.AnotherBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return instance;
    };
    $instance.AnotherThing = function () {
      var $this = this;
      return;
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

  this.$class('1659f51c', 'MainClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.maybe($g.asyncinheritance.FirstClass.new()).then(function (value) {
        instance.FirstClass = value;
      }));
      init.push($promise.maybe($g.asyncinheritance.SecondClass.new()).then(function (value) {
        instance.SecondClass = value;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function () {
      var $this = this;
      return $t.fastbox($this.SomeBool.$wrapped && !$this.AnotherBool.$wrapped, $g.____testlib.basictypes.Boolean);
    };
    Object.defineProperty($instance, 'SomeBool', {
      get: function () {
        return this.FirstClass.SomeBool;
      },
      set: function (val) {
        this.FirstClass.SomeBool = val;
      },
    });
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        return this.SecondClass.AnotherBool;
      },
      set: function (val) {
        this.SecondClass.AnotherBool = val;
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
        "DoSomething|2|29dc432d<43834c3f>": true,
        "AnotherThing|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('49e8850e', function () {
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.asyncinheritance.MainClass.new()).then(function ($result0) {
              $result = $result0.DoSomething();
              $current = 1;
              $continue($resolve, $reject);
              return;
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
  });
});
