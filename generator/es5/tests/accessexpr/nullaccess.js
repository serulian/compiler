$module('nullaccess', function () {
  var $static = this;
  this.$class('c3e4f083', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeBool = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeBool|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    var sc2;
    var sc3;
    sc = $g.nullaccess.SomeClass.new();
    sc2 = $g.nullaccess.SomeClass.new();
    sc3 = null;
    return $t.fastbox((sc.SomeBool().$wrapped && $t.syncnullcompare($t.dynamicaccess(sc2, 'SomeBool', false), function () {
      return $t.fastbox(false, $g.____testlib.basictypes.Boolean);
    }).$wrapped) && $t.syncnullcompare($t.dynamicaccess(sc3, 'SomeBool', false), function () {
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    }).$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
