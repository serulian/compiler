$module('basic', function () {
  var $static = this;
  this.$class('78afdb83', 'SomeClass', false, '', function () {
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
      return $t.fastbox(32, $g.____testlib.basictypes.Integer);
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
        "Declare|1|fd8bc7c9<78afdb83>": true,
        "GetValue|2|fd8bc7c9<bb8d3aad>": true,
        "GetMainValue|2|fd8bc7c9<bb8d3aad>": true,
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
        "GetValue|2|fd8bc7c9<bb8d3aad>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$agent('96b99731', 'SomeAgent', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.GetMainValue = function () {
      var $this = this;
      return $t.fastbox($this.$principal.GetValue().$wrapped + 10, $g.____testlib.basictypes.Integer);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetMainValue|2|fd8bc7c9<bb8d3aad>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.basic.SomeClass.Declare();
    return $t.fastbox(sc.GetMainValue().$wrapped == 42, $g.____testlib.basictypes.Boolean);
  };
});
