$module('unary', function () {
  var $static = this;
  this.$class('44d9119f', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $static.$not = function (instance) {
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.DoSomething = function (first) {
    $g.unary.SomeClass.$not(first);
    return;
  };
});
