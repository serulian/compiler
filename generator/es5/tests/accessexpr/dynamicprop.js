$module('dynamicprop', function () {
  var $static = this;
  this.$class('ab4079a0', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.value = $t.fastbox(42, $g.________testlib.basictypes.Integer);
      return instance;
    };
    $instance.set$SomeProp = function (val) {
      var $this = this;
      $this.value = val;
      return;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      return $this.value;
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProp|3|f22a98b6": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    var sca;
    sc = $g.dynamicprop.SomeClass.new();
    sc.set$SomeProp($t.fastbox(123, $g.________testlib.basictypes.Integer));
    sca = sc;
    return $t.fastbox(($t.cast($t.dynamicaccess(sca, 'SomeProp', false), $g.________testlib.basictypes.Integer, false).$wrapped == 123) && (sc.SomeProp().$wrapped == 123), $g.________testlib.basictypes.Boolean);
  };
});
