$module('interface', function () {
  var $static = this;
  this.$class('87868828', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      return $t.fastbox(42, $g.________testlib.basictypes.Integer);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeValue|3|db1c26c2": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('d317a9be', 'Valuable', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeValue|3|db1c26c2": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('f7b1981d', 'Valued', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.interface.Valuable;
    };
    $instance.GetValue = function () {
      var $this = this;
      return $this.$wrapped.SomeValue();
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetValue|2|fb1385bf<db1c26c2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    var v;
    sc = $g.interface.SomeClass.new();
    v = sc;
    return $t.fastbox($t.fastbox(v, $g.interface.Valued).GetValue().$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
