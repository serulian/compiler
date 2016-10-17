$module('nullaccess', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.SomeBool = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['SomeBool', 3, $g.____testlib.basictypes.Boolean.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.nullaccess.SomeClass).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var sc;
    var sc2;
    var sc3;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.nullaccess.SomeClass.new().then(function ($result0) {
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
            $g.nullaccess.SomeClass.new().then(function ($result0) {
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
            sc2 = $result;
            sc3 = null;
            sc.SomeBool().then(function ($result2) {
              return $promise.resolve($t.unbox($result2)).then(function ($result1) {
                return ($promise.shortcircuit($result1, true) || $t.dynamicaccess(sc2, 'SomeBool')).then(function ($result4) {
                  return $promise.resolve($result4).then(function ($result3) {
                    return $promise.resolve($result1 && $t.unbox($t.nullcompare($result3, $t.box(false, $g.____testlib.basictypes.Boolean)))).then(function ($result0) {
                      return ($promise.shortcircuit($result0, true) || $t.dynamicaccess(sc3, 'SomeBool')).then(function ($result6) {
                        return $promise.resolve($result6).then(function ($result5) {
                          $result = $t.box($result0 && $t.unbox($t.nullcompare($result5, $t.box(true, $g.____testlib.basictypes.Boolean))), $g.____testlib.basictypes.Boolean);
                          $current = 3;
                          $continue($resolve, $reject);
                          return;
                        });
                      });
                    });
                  });
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
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
