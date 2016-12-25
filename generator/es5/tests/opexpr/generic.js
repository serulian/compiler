$module('generic', function () {
  var $static = this;
  this.$class('63802a70', 'SomeClass', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.BoolValue = $t.property(function () {
      var $this = this;
      return $t.fastbox(false, $g.____testlib.basictypes.Boolean);
    });
    $static.$bool = function (value) {
      return value.BoolValue();
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "BoolValue|3|43834c3f": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.generic.SomeClass($g.____testlib.basictypes.Integer).new();
    return $t.fastbox(!$g.generic.SomeClass($g.____testlib.basictypes.Integer).$bool(sc).$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
