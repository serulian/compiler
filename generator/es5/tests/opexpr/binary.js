$module('binary', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
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
    this.$typesig = function () {
      var computed = $t.createtypesig(['xor', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['or', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['and', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['leftshift', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['not', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['plus', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['minus', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['times', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['div', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['mod', 4, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.binary.SomeClass).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  $static.DoSomething = function (first, second) {
    var $result;
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

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Integer.$plus($t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $g.____testlib.basictypes.Integer.$equals($result1, $t.box(3, $g.____testlib.basictypes.Integer)).then(function ($result0) {
                $result = $result0;
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
