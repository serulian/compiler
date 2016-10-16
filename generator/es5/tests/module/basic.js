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
    return $promise.new(function (resolve) {
      $static.someInt = $t.box(true, $g.____testlib.basictypes.Boolean);
      resolve();
    });
  }, 'af4b3683', []);
  this.$init(function () {
    return $g.basic.AnotherFunction().then(function ($result0) {
      $static.anotherBool = $result0;
    });
  }, '8795b5d7', ['af4b3683']);
});
