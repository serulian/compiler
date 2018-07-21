$module('basic', function () {
  var $static = this;
  this.$class('eaf2f2d0', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeAgent) {
      var instance = new $static();
      instance.SomeAgent = SomeAgent;
      instance.SomeAgent.$principal = instance;
      return instance;
    };
    $static.Declare = function () {
      return $g.basic.SomeClass.new($g.basic.SomeAgent.new());
    };
    $instance.GetValue = function () {
      var $this = this;
      return $t.fastbox(32, $g.________testlib.basictypes.Integer);
    };
    Object.defineProperty($instance, 'GetMainValue', {
      get: function () {
        return this.SomeAgent.GetMainValue.bind(this.SomeAgent);
      },
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Declare|1|cf412abd<eaf2f2d0>": true,
        "GetValue|2|cf412abd<2e508ae6>": true,
        "GetMainValue|2|cf412abd<2e508ae6>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('f752a75d', 'SomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetValue|2|cf412abd<2e508ae6>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$agent('a48d0eec', 'SomeAgent', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.GetMainValue = function () {
      var $this = this;
      return $t.fastbox($this.$principal.GetValue().$wrapped + 10, $g.________testlib.basictypes.Integer);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetMainValue|2|cf412abd<2e508ae6>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.basic.SomeClass.Declare();
    return $t.fastbox(sc.GetMainValue().$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
