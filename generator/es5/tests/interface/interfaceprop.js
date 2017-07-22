$module('interfaceprop', function () {
  var $static = this;
  this.$class('03095ee7', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.propValue = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      return instance;
    };
    $instance.set$SomeProperty = function (val) {
      var $this = this;
      $this.propValue = val;
      return;
    };
    $instance.SomeProperty = $t.property(function () {
      var $this = this;
      return $this.propValue;
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProperty|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('48a8f275', 'AnotherClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.set$SomeProperty = function (val) {
      var $this = this;
      return;
    };
    $instance.SomeProperty = $t.property($t.markpromising(function () {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $promise.translate($g.interfaceprop.DoSomethingAsync()).then(function ($result0) {
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
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    }));
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProperty|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('ec209969', 'SomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProperty|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('c342c45f', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var si;
    var si2;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            si = $g.interfaceprop.SomeClass.new();
            si2 = $g.interfaceprop.AnotherClass.new();
            $promise.maybe(si.set$SomeProperty($t.fastbox(false, $g.________testlib.basictypes.Boolean))).then(function ($result0) {
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
            $promise.maybe(si.SomeProperty()).then(function ($result1) {
              return $promise.resolve(!$result1.$wrapped).then(function ($result0) {
                return ($promise.shortcircuit($result0, true) || $promise.maybe(si2.SomeProperty())).then(function ($result2) {
                  $result = $t.fastbox($result0 && $result2.$wrapped, $g.________testlib.basictypes.Boolean);
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
  });
});
