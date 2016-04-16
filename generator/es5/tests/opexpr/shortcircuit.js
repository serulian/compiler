$module('shortcircuit', function () {
  var $static = this;
  this.$class('someError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($t.box('WHY CALLED? ', $g.____testlib.basictypes.String));
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    });
  });

  $static.neverCalled = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.shortcircuit.someError.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $state.reject($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  $static.anotherNeverCalled = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.shortcircuit.someError.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $state.reject($result);
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
            $promise.resolve(false).then(function ($result1) {
              return ($promise.shortcircuit($result1, false) || $g.shortcircuit.neverCalled()).then(function ($result2) {
                return $promise.resolve(!($result1 && $t.unbox($result2))).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || $g.shortcircuit.anotherNeverCalled()).then(function ($result3) {
                    $result = $t.box($result0 || $t.unbox($result3), $g.____testlib.basictypes.Boolean);
                    $state.current = 1;
                    $callback($state);
                  });
                });
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
