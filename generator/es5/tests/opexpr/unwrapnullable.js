$module('unwrapnullable', function () {
  var $static = this;
  this.$class('789ee0d6', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeProperty = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProperty|3|71258460": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('3753f9a6', 'SomeType', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.unwrapnullable.SomeClass;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var first;
    var second;
    var st;
    var st2;
    st = $t.fastbox($g.unwrapnullable.SomeClass.new(), $g.unwrapnullable.SomeType);
    st2 = null;
    first = $t.syncnullcompare($t.dynamicaccess($t.unbox(st), 'SomeProperty', false), function () {
      return $t.fastbox(false, $g.________testlib.basictypes.Boolean);
    });
    second = $t.syncnullcompare($t.dynamicaccess($t.unbox(st2), 'SomeProperty', false), function () {
      return $t.fastbox(false, $g.________testlib.basictypes.Boolean);
    });
    return $t.fastbox(first.$wrapped && !second.$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
