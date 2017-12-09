$module('classcastfail', function () {
  var $static = this;
  this.$class('aedde0c1', 'SomeClass', false, '', function () {
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

  this.$class('d2da6853', 'AnotherClass', false, '', function () {
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

  $static.TEST = function () {
    var ac;
    var sc;
    sc = $g.classcastfail.SomeClass.new();
    ac = $g.classcastfail.AnotherClass.new();
    $t.cast(ac, $g.classcastfail.SomeClass, false);
    return;
  };
});
