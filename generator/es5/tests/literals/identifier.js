$module('identifier', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.AnotherFunction = function () {
    return $promise.empty();
  };
  $static.DoSomething = function (someParam) {
    var someVar;
    var $state = $t.sm(function ($continue) {
      someVar = $t.box(2, $g.____testlib.basictypes.Integer);
      $g.identifier.SomeClass;
      $g.identifier.AnotherFunction;
      someParam;
      someVar;
      $state.resolve();
    });
    return $promise.build($state);
  };
});
