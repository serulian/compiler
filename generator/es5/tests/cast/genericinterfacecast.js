$module('genericinterfacecast', function () {
  var $static = this;
  this.$class('58a4c803', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeValue|3|43834c3f": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('3d5b9895', 'SomeInterface', true, '', function (T) {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
      };
      computed[("SomeValue|3|" + $t.typeid(T)) + ""] = true;
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.genericinterfacecast.SomeClass.new();
    return $t.cast(sc, $g.genericinterfacecast.SomeInterface($g.____testlib.basictypes.Boolean), false).SomeValue();
  };
});
