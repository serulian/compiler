$module('shortcircuit', function () {
  var $static = this;
  this.$class('someError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box('WHY CALLED? ', $g.____testlib.basictypes.String));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Message', 3, $g.____testlib.basictypes.String.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.shortcircuit.someError).$typeref()]);
    };
  });

  $static.neverCalled = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.shortcircuit.someError.new().then(function ($result0) {
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
            $reject($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.anotherNeverCalled = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.shortcircuit.someError.new().then(function ($result0) {
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
            $reject($result);
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
            $promise.resolve(false).then(function ($result1) {
              return ($promise.shortcircuit($result1, true) || $g.shortcircuit.neverCalled()).then(function ($result2) {
                return $promise.resolve(!($result1 && $t.unbox($result2))).then(function ($result0) {
                  return ($promise.shortcircuit($result0, false) || $g.shortcircuit.anotherNeverCalled()).then(function ($result3) {
                    $result = $t.box($result0 || $t.unbox($result3), $g.____testlib.basictypes.Boolean);
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                });
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
