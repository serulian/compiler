$module('compare', function () {
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
    $static.$equals = function (first, second) {
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
    $static.$compare = function (first, second) {
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($t.nominalwrap(1, $g.____testlib.basictypes.Integer));
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

  $static.TEST = function () {
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.compare.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            first = $result;
            $g.compare.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            second = $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $t.nominalwrap(!$t.nominalunwrap($result0), $g.____testlib.basictypes.Boolean);
              $state.current = 4;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.nominalwrap($t.nominalunwrap($result0) < 0, $g.____testlib.basictypes.Boolean);
              $state.current = 5;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.nominalwrap($t.nominalunwrap($result0) > 0, $g.____testlib.basictypes.Boolean);
              $state.current = 6;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 6:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.nominalwrap($t.nominalunwrap($result0) <= 0, $g.____testlib.basictypes.Boolean);
              $state.current = 7;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 7:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.nominalwrap($t.nominalunwrap($result0) >= 0, $g.____testlib.basictypes.Boolean);
              $state.current = 8;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 8:
            $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $result0;
              $state.current = 9;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 9:
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
