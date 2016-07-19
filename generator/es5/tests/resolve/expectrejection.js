$module('expectrejection', function () {
  var $static = this;
  this.$class('SimpleError', false, '', function () {
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
        $resolve($t.box('yo!', $g.____testlib.basictypes.String));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Message', 3, $g.____testlib.basictypes.String.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.expectrejection.SimpleError).$typeref()]);
    };
  });

  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.expectrejection.SimpleError.new().then(function ($result0) {
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
    var a;
    var b;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.expectrejection.DoSomething().then(function ($result0) {
              a = $result0;
              b = null;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function ($rejected) {
              b = $rejected;
              a = null;
              $current = 1;
              $continue($resolve, $reject);
              return;
            });
            return;

          case 1:
            $promise.resolve(a == null).then(function ($result0) {
              return ($promise.shortcircuit($result0, true) || $t.assertnotnull(b).Message()).then(function ($result2) {
                return ($promise.shortcircuit($result0, true) || $g.____testlib.basictypes.String.$equals($result2, $t.box('yo!', $g.____testlib.basictypes.String))).then(function ($result1) {
                  $result = $t.box($result0 && $t.unbox($result1), $g.____testlib.basictypes.Boolean);
                  $current = 2;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
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
