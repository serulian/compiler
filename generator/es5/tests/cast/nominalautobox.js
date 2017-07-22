$module('nominalautobox', function () {
  var $static = this;
  this.$class('f0465a3f', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$type('e9a6d7f9', 'SomeNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.nominalautobox.SomeClass;
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeValue|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.nominalautobox.SomeClass.new();
    return $t.cast(sc, $g.nominalautobox.SomeNominal, false).SomeValue();
  };
});
