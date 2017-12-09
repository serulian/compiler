$module('nominalcastfail', function () {
  var $static = this;
  this.$class('9a6729ae', 'SomeClass', false, '', function () {
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

  this.$class('8547cbac', 'AnotherClass', false, '', function () {
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

  this.$type('8fbb83dc', 'SomeNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.nominalcastfail.SomeClass;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var ac;
    var sc;
    sc = $g.nominalcastfail.SomeClass.new();
    ac = $g.nominalcastfail.AnotherClass.new();
    $t.cast(ac, $g.nominalcastfail.SomeNominal, false);
    return;
  };
});
