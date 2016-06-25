$module('calls', function () {
  var $static = this;
  this.$class('SomeError', false, '', function () {
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
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box('huh?', $g.____testlib.basictypes.String));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Message', 3, $g.____testlib.basictypes.String.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.calls.SomeError).$typeref()]);
    };
  });

  $static.getValue = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
  $static.failValue = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.calls.SomeError.new().then(function ($result0) {
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
  $static.getIntValue = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.box(45, $g.____testlib.basictypes.Integer));
      return;
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.calls.getIntValue().then(function ($result2) {
              return $g.____testlib.basictypes.Integer.$equals($result2, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || $g.calls.failValue()).then(function ($result3) {
                    return ($promise.shortcircuit($result0, false) || $g.calls.getValue()).then(function ($result4) {
                      $result = $result0 ? $result3 : $result4;
                      $current = 1;
                      $continue($resolve, $reject);
                      return;
                    });
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
