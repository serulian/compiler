$module('maybeasyncfunction', function () {
  var $static = this;
  this.$class('63735d1b', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SimpleFunction = $t.markpromising(function () {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        localasyncloop: while (true) {
          switch ($current) {
            case 0:
              $promise.translate($g.maybeasyncfunction.DoSomethingAsync()).then(function ($result0) {
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
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SimpleFunction|2|fb1385bf<71258460>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('93608c9a', 'AnotherClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SimpleFunction = function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SimpleFunction|2|fb1385bf<71258460>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('1f1ff333', 'ISimple', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SimpleFunction|2|fb1385bf<71258460>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('9cca6f95', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var isi;
    var isi2;
    var r1;
    var r2;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            isi = $g.maybeasyncfunction.SomeClass.new();
            isi2 = $g.maybeasyncfunction.AnotherClass.new();
            $promise.maybe(isi.SimpleFunction()).then(function ($result0) {
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
            r1 = $result;
            $promise.maybe(isi2.SimpleFunction()).then(function ($result0) {
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
            r2 = $result;
            $resolve($t.fastbox(r1.$wrapped && r2.$wrapped, $g.________testlib.basictypes.Boolean));
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
