$module('identifier', function () {
  var $static = this;
  this.$class('1de9520f', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.AnotherFunction = function () {
    return $promise.empty();
  };
  $static.DoSomething = function (someParam) {
    var someVar;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      someVar = $t.box(2, $g.____testlib.basictypes.Integer);
      $g.identifier.SomeClass;
      $g.identifier.AnotherFunction;
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
