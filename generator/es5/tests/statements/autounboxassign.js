$module('autounboxassign', function () {
  var $static = this;
  this.$class('3234898a', 'SomeClass', false, '', function () {
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
        "SomeValue|3|bb8d3aad": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('2cec5b57', 'AnotherType', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.autounboxassign.SomeClass;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var at;
    var sc;
    at = $t.fastbox($g.autounboxassign.SomeClass.new(), $g.autounboxassign.AnotherType);
    sc = null;
    sc = $t.unbox(at);
    return $t.fastbox($t.syncnullcompare($t.dynamicaccess(sc, 'SomeValue', false), function () {
      return $t.fastbox(0, $g.________testlib.basictypes.Integer);
    }).$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
