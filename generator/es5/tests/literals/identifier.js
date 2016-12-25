$module('identifier', function () {
  var $static = this;
  this.$class('1de9520f', 'SomeClass', false, '', function () {
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

  $static.AnotherFunction = function () {
    return;
  };
  $static.DoSomething = function (someParam) {
    var someVar;
    someVar = $t.fastbox(2, $g.____testlib.basictypes.Integer);
    $g.identifier.SomeClass;
    $g.identifier.AnotherFunction;
    return;
  };
});
