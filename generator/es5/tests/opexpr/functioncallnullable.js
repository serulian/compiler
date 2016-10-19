$module('functioncallnullable', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.SomeMethod = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      var computed = $t.createtypesig(['SomeMethod', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.functioncallnullable.SomeClass).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  this.$class('AnotherClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.AnotherMethod = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(false, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      var computed = $t.createtypesig(['AnotherMethod', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.functioncallnullable.AnotherClass).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var ac;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.functioncallnullable.SomeClass.new().then(function ($result0) {
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
            sc = $result;
            ac = null;
            $t.nullableinvoke(sc, 'SomeMethod', true, []).then(function ($result2) {
              return $promise.resolve($result2).then(function ($result1) {
                return $promise.resolve($t.unbox($t.nullcompare($result1, $t.box(false, $g.____testlib.basictypes.Boolean)))).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || $t.nullableinvoke(ac, 'AnotherMethod', true, [])).then(function ($result4) {
                    return $promise.resolve($result4).then(function ($result3) {
                      $result = $t.box($result0 && $t.unbox($t.nullcompare($result3, $t.box(true, $g.____testlib.basictypes.Boolean))), $g.____testlib.basictypes.Boolean);
                      $current = 2;
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
