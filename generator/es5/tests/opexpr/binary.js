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
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$or = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$and = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$leftshift = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$not = function (value) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(value);
        return;
      };
      return $promise.new($continue);
    };
    $static.$plus = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$minus = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(right);
        return;
      };
      return $promise.new($continue);
    };
    $static.$times = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$div = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
    $static.$mod = function (left, right) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(left);
        return;
      };
      return $promise.new($continue);
    };
  });

  $static.DoSomething = function (first, second) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.binary.SomeClass.$plus(first, second).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $result;
            $g.binary.SomeClass.$minus(first, second).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $result;
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Integer.$plus($t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.box(3, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
