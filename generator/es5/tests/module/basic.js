$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.resolve($g.basic.someInt);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.resolve($g.basic.anotherBool);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  this.$init(function () {
    return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
      $static.someInt = result;
    });
  });
  this.$init(function () {
    return $g.basic.AnotherFunction().then(function ($result0) {
      return $promise.resolve($result0);
    }).then(function (result) {
      $static.anotherBool = result;
    });
  });
});
