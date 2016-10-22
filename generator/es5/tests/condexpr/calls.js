$module('calls', function () {
  var $static = this;
  this.$class('165552b1', 'SomeError', false, '', function () {
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
        $resolve($t.fastbox('huh?', $g.____testlib.basictypes.String));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Message|3|538656f2": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.getValue = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
  $static.failValue = function () {
    var $result;
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
      $resolve($t.fastbox(45, $g.____testlib.basictypes.Integer));
      return;
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
            $g.calls.getIntValue().then(function ($result2) {
              return $g.____testlib.basictypes.Integer.$equals($result2, $t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                return $promise.resolve($result1.$wrapped).then(function ($result0) {
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
