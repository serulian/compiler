$module('structnew', function () {
  var $static = this;
  this.$class('cd38ba2a', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.SomeField = SomeField;
      instance.anotherField = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return instance;
    };
    $instance.set$AnotherField = function (val) {
      var $this = this;
      $this.anotherField = val;
      return;
    };
    $instance.AnotherField = $t.property(function () {
      var $this = this;
      return $this.anotherField;
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherField|3|f7f23c49": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $temp0;
    var sc;
    sc = ($temp0 = $g.structnew.SomeClass.new($t.fastbox(2, $g.____testlib.basictypes.Integer)), $temp0.set$AnotherField($t.fastbox(true, $g.____testlib.basictypes.Boolean)), $temp0);
    return $t.fastbox((sc.SomeField.$wrapped == 2) && sc.anotherField.$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
