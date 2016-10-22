$module('functioncallnullable', function () {
  var $static = this;
  this.$class('fa2a0786', 'SomeClass', false, '', function () {
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
        $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeMethod|2|29dc432d<5ab5941e>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('97cf7b5c', 'AnotherClass', false, '', function () {
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
        $resolve($t.fastbox(false, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherMethod|2|29dc432d<5ab5941e>": true,
      };
      return this.$cachedtypesig = computed;
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
                return $promise.resolve($t.nullcompare($result1, $t.fastbox(false, $g.____testlib.basictypes.Boolean)).$wrapped).then(function ($result0) {
                  return ($promise.shortcircuit($result0, true) || $t.nullableinvoke(ac, 'AnotherMethod', true, [])).then(function ($result4) {
                    return $promise.resolve($result4).then(function ($result3) {
                      $result = $t.fastbox($result0 && $t.nullcompare($result3, $t.fastbox(true, $g.____testlib.basictypes.Boolean)).$wrapped, $g.____testlib.basictypes.Boolean);
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
