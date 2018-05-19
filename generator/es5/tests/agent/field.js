$module('field', function () {
  var $static = this;
  this.$class('6c0ca75a', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeAgent) {
      var instance = new $static();
      instance.SomeAgent = SomeAgent;
      instance.SomeAgent.$principal = instance;
      return instance;
    };
    $static.Declare = function () {
      return $g.field.SomeClass.new($g.field.SomeAgent.new($t.fastbox(42, $g.________testlib.basictypes.Integer)));
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        return this.SomeAgent.SomeField;
      },
      set: function (val) {
        this.SomeAgent.SomeField = val;
      },
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Declare|1|fb1385bf<6c0ca75a>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('9a000095', 'SomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$agent('62ce5362', 'SomeAgent', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.SomeField = SomeField;
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.field.SomeClass.Declare();
    return $t.fastbox(sc.SomeField.$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
