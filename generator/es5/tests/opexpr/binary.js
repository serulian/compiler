$module('binary', function () {
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
    $static.$xor = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$or = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$and = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$leftshift = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$not = function (value) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(value);
        return;
      });
      return $promise.build($state);
    };
    $static.$plus = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$minus = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(right);
        return;
      });
      return $promise.build($state);
    };
    $static.$times = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$div = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
    $static.$mod = function (left, right) {
      var $state = $t.sm(function ($continue) {
        $state.resolve(left);
        return;
      });
      return $promise.build($state);
    };
  });

  $static.DoSomething = function (first, second) {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.binary.SomeClass.$plus(first, second).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
            $g.binary.SomeClass.$minus(first, second).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            $state.current = -1;
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
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.Integer.$plus($t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.box(3, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $state.current = 1;
                $continue($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $state.resolve($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
