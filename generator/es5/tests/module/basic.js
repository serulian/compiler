$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($g.basic.someInt);
      return;
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($g.basic.anotherBool);
      return;
    };
    return $promise.new($continue);
  };
  this.$init(function () {
    return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
      $static.someInt = result;
    });
  });
  this.$init(function () {
    return $g.basic.AnotherFunction().then(function ($result0) {
      $static.anotherBool = $result0;
    });
  });
});
