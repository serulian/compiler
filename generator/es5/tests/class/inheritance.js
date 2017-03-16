$module('inheritance', function () {
  var $static = this;
  this.$class('16519abc', 'FirstClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeBool = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      return instance;
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
        "DoSomething|2|89b8f38e<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('bcec6d2c', 'SecondClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
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
        "AnotherThing|2|89b8f38e<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('131aac76', 'MainClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.FirstClass = $g.inheritance.FirstClass.new();
      instance.SecondClass = $g.inheritance.SecondClass.new();
      return instance;
    };
    $instance.DoSomething = function () {
      var $this = this;
      return $this.SomeBool;
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
        "DoSomething|2|89b8f38e<f7f23c49>": true,
        "AnotherThing|2|89b8f38e<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    return $g.inheritance.MainClass.new().DoSomething();
  };
});
