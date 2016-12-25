$module('binary', function () {
  var $static = this;
  this.$class('e724e586', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $static.$xor = function (left, right) {
      return left;
    };
    $static.$or = function (left, right) {
      return left;
    };
    $static.$and = function (left, right) {
      return left;
    };
    $static.$leftshift = function (left, right) {
      return left;
    };
    $static.$not = function (value) {
      return value;
    };
    $static.$plus = function (left, right) {
      return left;
    };
    $static.$minus = function (left, right) {
      return right;
    };
    $static.$times = function (left, right) {
      return left;
    };
    $static.$div = function (left, right) {
      return left;
    };
    $static.$mod = function (left, right) {
      return left;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.DoSomething = function (first, second) {
    $g.binary.SomeClass.$plus(first, second);
    $g.binary.SomeClass.$minus(first, second);
    return;
  };
  $static.TEST = function () {
    return $t.fastbox((1 + 2) == 3, $g.____testlib.basictypes.Boolean);
  };
});
