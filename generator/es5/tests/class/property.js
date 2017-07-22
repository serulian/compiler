$module('property', function () {
  var $static = this;
  this.$class('88e36ba3', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
      return instance;
    };
    $instance.set$SomeProp = function (val) {
      var $this = this;
      $this.SomeBool = val;
      return;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      return $this.SomeBool;
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProp|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.AnotherFunction = function (sc) {
    sc.SomeProp();
    sc.set$SomeProp($t.fastbox(true, $g.________testlib.basictypes.Boolean));
    return sc.SomeProp();
  };
  $static.TEST = function () {
    return $g.property.AnotherFunction($g.property.SomeClass.new());
  };
});
